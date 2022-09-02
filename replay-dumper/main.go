package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/maxsupermanhd/go-wz/wznet"
)

var (
	filepath       = flag.String("f", "./replay.wzrp", "Path to replay to dump")
	statsdir       = flag.String("stats", "./data/mp/stats/", "Path to stats directory")
	mapout         = flag.String("mapout", "./map.wz", "Path to save embedded map. Use - to disable")
	short          = flag.Bool("short", false, "Do not print out everything")
	dOrder         = flag.Bool("dorder", false, "Dump unit commands")
	dResearch      = flag.Bool("dres", true, "Dump research packets")
	dStructinfo    = flag.Bool("dstructs", false, "Dump structure orders (research/production)")
	checkIllegals  = flag.Bool("verifyUnits", true, "Verify units composition")
	netPlayPlayers = []NetplayPlayers{}
	namePadLength  = 2
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	flag.Parse()
	PrintNShort("Replay dumper starting up...")

	if *dStructinfo || *checkIllegals {
		log.Printf("Loading research from [%s]...", *statsdir)
		must(loadStatsData(*statsdir))
	}

	log.Printf("Dumping replay file file [%s]...", *filepath)
	f := noerr(os.Open(*filepath))

	PrintNShort("Reading magic...")
	readMagic(f)

	PrintNShort("Reading header JSON...")
	readSettings(f)

	PrintNShort("Reading embedded map data...")
	readEmbeddedMap(f)

	log.Printf("Begin of netmessages...")
	readNetMessages(f)

	PrintNShort("Bye!")
	f.Close()
}

func readNetMessages(f *os.File) {
	total := 0
	gameTime := uint32(0)

	netsizes := map[string]uint64{}

	playersautorepair := map[string]uint32{}

	for {
		pPlayer, _ := noerr2(wznet.ReadByte(f))
		pType, _ := noerr2(wznet.ReadByte(f))
		l := uint32(0)
		success := false
		b := byte(0)
		n := uint(0)
		for {
			b, success = noerr2(wznet.ReadByte(f))
			r := false
			r, l = wznet.Decode_uint32_t(b, l, n)
			n++
			if r {
				break
			}
		}
		if !success {
			log.Printf("Not success, len %d", l)
			break
		}
		if pType == wznet.REPLAY_ENDED {
			log.Printf("End of net messages. (%d total)", total)
			break
		}
		data := noerr(wznet.ReadBytes(f, int(l)))
		msgid, ok := wznet.NetMessageType[pType]
		size1, ok1 := netsizes[msgid]
		if ok1 {
			netsizes[msgid] = size1 + uint64(l)
		} else {
			netsizes[msgid] = uint64(l)
		}
		rprint := func(f string, v ...interface{}) {
			vv := []interface{}{gameTime, gameTimeToString(gameTime), -namePadLength, netPlayPlayers[pPlayer].Name}
			vv = append(vv, v...)
			log.Printf("(%7d % 9s) [% *s] "+f, vv...)
		}
		if !ok {
			log.Printf("Unknown message type: %d", pType)
		} else {
			r := bytes.NewReader(data)
			// log.Printf("Message from %2d %-28s len %3d %3d", pPlayer, msgid, len(data), l)
			switch pType {
			case wznet.GAME_GAME_TIME:
				_ = noerr(wznet.NETreadU32(r)) // latencyTicks :
				gameTime = noerr(wznet.NETreadU32(r))
				_ = noerr(wznet.NETreadU16(r)) // crc :
				_ = noerr(wznet.NETreadU16(r)) // wantedLatency :
			case wznet.GAME_DEBUG_MODE:
				enable := noerr(wznet.NETreadU8(r))
				rprint("Debug mode request %v", enable)
			case wznet.GAME_DEBUG_ADD_STRUCTURE:
				fallthrough
			case wznet.GAME_DEBUG_ADD_DROID:
				fallthrough
			case wznet.GAME_DEBUG_ADD_FEATURE:
				fallthrough
			case wznet.GAME_DEBUG_REMOVE_DROID:
				fallthrough
			case wznet.GAME_DEBUG_REMOVE_FEATURE:
				fallthrough
			case wznet.GAME_DEBUG_REMOVE_STRUCTURE:
				fallthrough
			case wznet.GAME_DEBUG_FINISH_RESEARCH:
				rprint("DEBUG %s", msgid)
			case wznet.GAME_STRUCTUREINFO:
				if !*checkIllegals && !*dStructinfo {
					break
				}
				_ = noerr(wznet.NETreadU8(r)) // player
				// if player != pPlayer {
				// 	rprint("Wrong player packet %v", player)
				// }
				_ = noerr(wznet.NETreadU32(r)) // structID
				structInfo := noerr(wznet.NETreadU8(r))
				if structInfo == 0 {
					droid := DroidDef{
						noerr(wznet.NETstring(r)),  // droid Name
						noerr(wznet.NETreadU32(r)), // droid ID
						noerr(wznet.NETreadS32(r)), // droid Type
						noerr(wznet.NETreadU8(r)),  // droid Body
						noerr(wznet.NETreadU8(r)),  // droid Brain
						noerr(wznet.NETreadU8(r)),  // droid Propulsion
						noerr(wznet.NETreadU8(r)),  // droid Repairunit
						noerr(wznet.NETreadU8(r)),  // droid Ecm
						noerr(wznet.NETreadU8(r)),  // droid Sensor
						noerr(wznet.NETreadU8(r)),  // droid Construct
					}
					droidNumWeapons := noerr(wznet.NETreadU8(r))
					droidWeapons := []uint32{}
					for i := uint8(0); i < droidNumWeapons; i++ {
						droidWeapons = append(droidWeapons, noerr(wznet.NETreadU32(r)))
					}
					if len(droidWeapons) > 2 {
						rprint("Built droid with more than 2 turrets! %v", droid)
					}
					droidTyped := typifyDroid(droid)
					if droidTyped != nil {
						if *dStructinfo {
							rprint("droid: %v", spew.Sdump(*droidTyped))
						}
						if illegalmsg := checkDroidIllegal(*droidTyped); illegalmsg != "" {
							rprint("Illegal droid (%v): %v", illegalmsg, spew.Sdump(*droidTyped))
						}
						if droidTyped.Repairunit == "AutoRepair" {
							_, ok := playersautorepair[netPlayPlayers[pPlayer].Name]
							if !ok {
								playersautorepair[netPlayPlayers[pPlayer].Name] = gameTime
							}
						}
					}
				}
			case wznet.GAME_DROIDINFO:
				if !*dOrder {
					break
				}
				player := noerr(wznet.NETreadU8(r))
				if player != pPlayer {
					log.Printf("Player missmatch in %s (%d netmessage %d packet)", msgid, pPlayer, player)
				}
				subtype := wznet.DroidOrderSybType(noerr(wznet.NETreadU32(r)))
				printparams := []interface{}{subtype.String()}
				switch subtype {
				case wznet.DroidOrderSybTypeObj:
					fallthrough
				case wznet.DroidOrderSybTypeLoc:
					order := wznet.DORDER(noerr(wznet.NETreadU32(r)))
					printparams = append(printparams, order.String())
					if subtype == wznet.DroidOrderSybTypeObj {
						destID := noerr(wznet.NETreadU32(r))
						destType := noerr(wznet.NETreadU32(r))
						printparams = append(printparams, "destID", destID, "destType", destType)
					} else {
						coordx := noerr(wznet.NETreadS32(r))
						coordy := noerr(wznet.NETreadS32(r))
						printparams = append(printparams, "x", coordx, "y", coordy, "tile aligned", coordx%128 == 0, coordy%128 == 0)
					}
					if order == wznet.DORDER_BUILD || order == wznet.DORDER_LINEBUILD {
						structref := noerr(wznet.NETreadU32(r))
						direction := noerr(wznet.NETreadU16(r))
						printparams = append(printparams, "structref", structref, "direction", direction)
					}
					if order == wznet.DORDER_LINEBUILD {
						coordx2 := noerr(wznet.NETreadS32(r))
						coordy2 := noerr(wznet.NETreadS32(r))
						printparams = append(printparams, "x2", coordx2, "y2", coordy2)
					}
					if order == wznet.DORDER_BUILDMODULE {
						index := noerr(wznet.NETreadU32(r))
						printparams = append(printparams, "index", index)
					}
					add := noerr(wznet.NETreadU8(r))
					printparams = append(printparams, "add", add)
				case wznet.DroidOrderSybTypeSec:
					secorder := noerr(wznet.NETreadU32(r))
					secstate := noerr(wznet.NETreadU32(r))
					printparams = append(printparams, "secorder", secorder, "secstate", secstate)
				}
				num := noerr(wznet.NETreadU32(r))
				droidids := []uint32{}
				droiddelta := uint32(0)
				for i := uint32(0); i < num; i++ {
					droiddelta += noerr(wznet.NETreadU32(r))
					droidids = append(droidids, droiddelta)
				}
				printparams = append(printparams, droidids)
				rprint(strings.Repeat("%v ", len(printparams)), printparams...)
			case wznet.GAME_RESEARCHSTATUS:
				if !*dResearch {
					break
				}
				player := noerr(wznet.NETreadU8(r))
				start := noerr(wznet.NETreadU8(r))
				building := noerr(wznet.NETreadU32(r))
				_ = building
				topic := noerr(wznet.NETreadU32(r))
				topicname := fmt.Sprint(topic)
				_ = topicname
				if int(topic) < len(aResearch) {
					topicname = aResearch[topic].Name
				} else {
					log.Printf("Topic overflow or underflow, topic %d total %d", topic, len(aResearch))
				}
				if player != pPlayer {
					log.Printf("Player missmatch in %s (%d netmessage %d packet)", msgid, pPlayer, player)
				}
				action := "picked "
				if start == 0 {
					action = "dropped"
				}
				_ = action
				if !*short {
					rprint("Research topic %s %s on building %d", topicname, action, building)
				}
				if r.Len() != 0 {
					spew.Dump(data)
					log.Printf("Did not parsed %d bytes", r.Len())
				}
			}
		}
		total++
	}
	log.Printf("Replay time: %v (%v ticks)", gameTimeToString(gameTime), gameTime)
	if !*short {
		log.Println("Replay packets sizes (bytes):")
		for msg, size := range netsizes {
			log.Printf("\t%v: %v", msg, size)
		}
		log.Println("Players auto-repair unit ecm (ticks):")
		for pname, gtime := range playersautorepair {
			log.Printf("\t%v: %v", pname, gameTimeToString(gtime))
		}
	}
}

func readEmbeddedMap(f *os.File) {
	dv := noerr(wznet.ReadUBE32(f))
	if dv != 1 {
		log.Printf("Embedded map data version is odd: %d, should be 1", dv)
	}
	len := noerr(wznet.ReadUBE32(f))
	if len != 0 {
		PrintNShort("Embedded map data len %d", len)
		b := noerr(wznet.ReadBytes(f, int(len)))
		if *mapout != "-" {
			must(os.WriteFile(*mapout, b, 0644))
		}
	} else {
		PrintNShort("Embedded map data is empty")
	}
}

func readSettings(f *os.File) {
	b := noerr(wznet.ReadBytes(f, int(noerr(wznet.ReadUBE32(f)))))
	var s ReplaySettings
	must(json.Unmarshal(b, &s))
	if s.ReplayFormatVer != 2 {
		log.Printf("Replay format version is odd: %d, should be 2", s.ReplayFormatVer)
	}
	log.Printf("Replay version %s (netcode %d.%d)", s.GameOptions.VersionString, s.Major, s.Minor)
	log.Println("Players:")
	for i, p := range s.GameOptions.NetplayPlayers {
		if p.Allocated {
			id := sha256.Sum256(noerr(base64.StdEncoding.DecodeString(s.GameOptions.Multistats[i].Identity)))
			log.Printf("Position %2d index %2d name [%s] [%s] [%x]", p.Position, i, p.Name, s.GameOptions.Multistats[i].Identity, id)
			if namePadLength < len(p.Name) {
				namePadLength = len(p.Name)
			}
		}
	}
	netPlayPlayers = s.GameOptions.NetplayPlayers
}

func readMagic(f *os.File) {
	magic := make([]byte, 4)
	if noerr(f.Read(magic)) != 4 {
		log.Fatal("Read failed to read magic")
	}
	if !bytes.Equal(magic, []byte{'W', 'Z', 'r', 'p'}) {
		log.Fatalf("Magic is not `WZRP`! Got [%4s]", magic)
	}
	PrintNShort("Magic is valid")
}

func gameTimeToString(gt uint32) string {
	return (time.Duration(int(gt/1000)) * time.Second).String()
}

func PrintNShort(format string, args ...interface{}) {
	if !*short {
		log.Printf(format, args...)
	}
}
