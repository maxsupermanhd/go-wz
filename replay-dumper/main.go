package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/image/draw"

	"github.com/davecgh/go-spew/spew"
	"github.com/dustin/go-heatmap"
	"github.com/dustin/go-heatmap/schemes"
	"github.com/dustin/go-humanize"
	"github.com/maxsupermanhd/go-wz/phobos"
	"github.com/maxsupermanhd/go-wz/wznet"
)

var (
	filepath             = flag.String("f", "./replay.wzrp", "Path to replay to dump, can be a url if fetch is true")
	fetch                = flag.Bool("fetch", false, "If true treat filepath as url and fetch replay into memory from it")
	statsdir             = flag.String("stats", "./data/mp/stats/", "Path to stats directory")
	mapout               = flag.String("mapout", "./map.wz", "Path to save embedded map. Use - to disable")
	short                = flag.Bool("short", false, "Do not print out everything")
	dOrder               = flag.Bool("dorder", false, "Dump unit commands")
	dResearch            = flag.Bool("dres", true, "Dump research packets")
	dStructinfo          = flag.Bool("dstructs", false, "Dump structure orders (research/production)")
	dumpUnits            = flag.Bool("dumpUnits", false, "Dump production units")
	checkIllegals        = flag.Bool("verifyUnits", true, "Verify units composition")
	genHeatmap           = flag.Bool("genHeatmap", false, "Generate heatmap of clicks")
	genHeatmapPath       = flag.String("heatmapOut", "./heatmap.png", "Path out for heatmap")
	heatmapIntensity     = flag.Int("heatmapIntensity", 20, "The impact size of each point on the output")
	heatmapScale         = flag.Int("mapZ", 32, "Scale of heatmap")
	mapW                 = flag.Int("mapW", 2048, "Map width")
	mapH                 = flag.Int("mapH", 2048, "Map height")
	phobosInfo           = flag.Bool("phobosInfo", true, "Fetch information about map from wz2100.euphobos.net/maps/")
	phobosMapSizeFromScr = flag.Bool("phobosMapSizeFromScr", false, "Get map size from map_scr or map_size")
	phobosPreview        = flag.Bool("phobosPreview", true, "Fetch map preview from wz2100.euphobos.net/maps/")
	filePreview          = flag.String("filePreview", "", "Overlay map preview (png) from this path")
	overrideMapHash      = flag.String("overrideHash", "", "Optional override for map hash")
	netPlayPlayers       = []NetplayPlayers{}
	namePadLength        = 2
	// clickHeatmap     = map[int]clickPoint{}
	clickHeatmap  = []heatmap.DataPoint{}
	mapHash       = ""
	replayOptions = GameOptions{}
)

// type clickPoint struct {
// 	Pos  heatmap.DataPoint
// 	Tick uint32
// }

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	flag.Parse()
	PrintNShort("Replay dumper starting up...")

	if *dResearch || *dStructinfo || *checkIllegals {
		log.Printf("Loading stats from [%s]...", *statsdir)
		must(loadStatsData(*statsdir))
	}

	var f *bytes.Buffer
	if *fetch {
		log.Printf("Fetching replay file [%s]...", *filepath)
		f = bytes.NewBuffer(noerr(io.ReadAll(noerr(http.Get(*filepath)).Body)))
	} else {
		log.Printf("Dumping replay file [%s]...", *filepath)
		f = bytes.NewBuffer(noerr(os.ReadFile(*filepath)))
	}

	PrintNShort("Reading magic...")
	readMagic(f)

	PrintNShort("Reading header JSON...")
	readSettings(f)

	PrintNShort("Reading embedded map data...")
	readEmbeddedMap(f)

	log.Printf("Begin of netmessages...")
	readNetMessages(f)

	if *genHeatmap {
		mw := *mapW
		mh := *mapH
		if *phobosInfo {
			log.Println("Fetching info about map...")
			info := noerr(phobos.FetchOnePhobosInfo(mapHash))
			if *phobosMapSizeFromScr {
				mw = info.MapScrX2
				mh = info.MapScrY2
			} else {
				fmt.Sscanf(info.MapSize, "%dx%d", &mw, &mh)
			}
			log.Printf("Size: W %d H %d", mw, mh)
		}
		log.Println("Generating heatmap...")
		clickHeatmap = append(clickHeatmap, heatmap.P(0, 0))
		clickHeatmap = append(clickHeatmap, heatmap.P(float64(mw*128), float64(mh*128)))
		for i, v := range clickHeatmap {
			clickHeatmap[i] = heatmap.P(v.X(), v.Y()*-1+float64(mh*128))
			if int(clickHeatmap[i].X()) > mw*128 {
				clickHeatmap[i] = heatmap.P(float64(mw*128), clickHeatmap[i].Y())
			}
			if int(clickHeatmap[i].Y()) > mh*128 {
				clickHeatmap[i] = heatmap.P(clickHeatmap[i].X(), float64(mw*128))
			}
		}
		log.Printf("Plotting %d samples...", len(clickHeatmap))
		hm := heatmap.Heatmap(image.Rect(0, 0, *heatmapScale*mw, *heatmapScale*mh), clickHeatmap, *heatmapIntensity, 255, schemes.AlphaFire)
		if *filePreview != "" {
			b := noerr(os.ReadFile(*filePreview))
			prv := noerr(png.Decode(bytes.NewBuffer(b)))
			log.Println("Rendering heightmap with preview...")
			i := image.NewRGBA(image.Rect(0, 0, *heatmapScale*mw, *heatmapScale*mh))
			draw.NearestNeighbor.Scale(i, i.Rect, prv, prv.Bounds(), draw.Over, nil)
			draw.Draw(i, i.Bounds(), hm, image.Point{}, draw.Over)
			hm = i
		} else if *phobosPreview {
			log.Println("Fetching map preview...")
			prv := noerr(phobos.FetchMapPreview(mapHash, phobos.PreviewTypeLargeJPEG))
			log.Println("Rendering heightmap with preview...")
			i := image.NewRGBA(image.Rect(0, 0, *heatmapScale*mw, *heatmapScale*mh))
			draw.NearestNeighbor.Scale(i, i.Rect, prv, prv.Bounds(), draw.Over, nil)
			draw.Draw(i, i.Bounds(), hm, image.Point{}, draw.Over)
			hm = i
		}
		log.Printf("Encoding heatmap to %q...", *genHeatmapPath)
		b := bytes.NewBuffer([]byte{})
		must(png.Encode(b, hm))
		must(os.WriteFile(*genHeatmapPath, b.Bytes(), 0644))
	}

	PrintNShort("Bye!")
}

func readNetMessages(f *bytes.Buffer) {
	total := 0
	gameTime := uint32(0)

	netsizes := map[string]uint64{}
	netcounts := map[string]uint64{}

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
		counts1, ok1 := netcounts[msgid]
		if ok1 {
			netcounts[msgid] = counts1 + 1
		} else {
			netcounts[msgid] = 1
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
				_ = noerr(wznet.NETreadU8(r)) // player
				// if player != pPlayer {
				// 	log.Printf("Player missmatch in %s (%d netmessage %d packet)", msgid, pPlayer, player)
				// }
				structID := noerr(wznet.NETreadU32(r))
				structInfo := noerr(wznet.NETreadU8(r))
				printparams := []interface{}{}
				if *dStructinfo {
					printparams = append(printparams, "structid", structID)
					printparams = append(printparams, "structinfo", structInfo)
				}
				if structInfo == wznet.STRUCTUREINFO_MANUFACTURE {
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
						if *dumpUnits {
							printparams = append(printparams, "droid", spew.Sdump(*droidTyped))
						}
						if *checkIllegals {
							if illegalmsg := checkDroidIllegal(*droidTyped); illegalmsg != "" {
								rprint("Illegal droid (%v): %v", illegalmsg, spew.Sdump(*droidTyped))
							}
						}
						if droidTyped.Repairunit == "AutoRepair" {
							_, ok := playersautorepair[netPlayPlayers[pPlayer].Name]
							if !ok {
								playersautorepair[netPlayPlayers[pPlayer].Name] = gameTime
							}
						}
					} else {
						if *dumpUnits {
							printparams = append(printparams, "droid UNTYPED", spew.Sdump(droid))
						}
					}
				}
				if len(printparams) > 0 {
					rprint(strings.Repeat("%v ", len(printparams)), printparams...)
				}
			case wznet.GAME_DROIDINFO:
				if !(*dOrder || *genHeatmap) {
					break
				}
				_ = noerr(wznet.NETreadU8(r))
				// if player != pPlayer {
				// 	log.Printf("Player missmatch in %s (%d netmessage %d packet)", msgid, pPlayer, player)
				// }
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
						if *genHeatmap {
							// clickHeatmap[replayOptions.NetplayPlayers[player].Position] = clickPoint{
							// 	Pos:  heatmap.P(float64(coordx), float64(coordy)),
							// 	Tick: gameTime,
							// }
							clickHeatmap = append(clickHeatmap, heatmap.P(float64(coordx), float64(coordy)))
						}
						printparams = append(printparams, "x", coordx, "y", coordy, "tile aligned", coordx%128 == 0, coordy%128 == 0)
					}
					if order == wznet.DORDER_BUILD || order == wznet.DORDER_LINEBUILD {
						structref := noerr(wznet.NETreadU32(r))
						direction := noerr(wznet.NETreadU16(r))
						printparams = append(printparams, "structref", refToStructName(structref), "direction", direction)
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
				oglen := len(droidids)
				printparams = append(printparams, "droids")
				if oglen > 4 {
					droidids = droidids[:4]
					printparams = append(printparams, droidids)
					printparams = append(printparams, fmt.Sprintf("and %d more", oglen-4))
				} else {
					printparams = append(printparams, droidids)
				}
				if *dOrder {
					rprint(strings.Repeat("%v ", len(printparams)), printparams...)
				}
			case wznet.GAME_RESEARCHSTATUS:
				if !*dResearch {
					break
				}
				_ = noerr(wznet.NETreadU8(r)) //player
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
				// if player != pPlayer {
				// 	log.Printf("Player missmatch in %s (%d netmessage %d packet)", msgid, pPlayer, player)
				// }
				action := "picked"
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
			case wznet.GAME_PLAYER_LEFT:
				rprint("Left.")
			case wznet.GAME_GIFT:
				// TODO parse gift
			default:
				rprint("%s not handled", msgid)
			}
		}
		total++
	}
	log.Printf("Replay time: %v (%v ticks)", gameTimeToString(gameTime), gameTime)
	if !*short {
		log.Println("Replay packets (bytes) (count):")
		for msg, size := range netsizes {
			log.Printf("\t%v: %v %v", msg, humanize.Bytes(size), netcounts[msg])
		}
		log.Println("Players auto-repair unit ecm (ticks):")
		for pname, gtime := range playersautorepair {
			log.Printf("\t%v: %v", pname, gameTimeToString(gtime))
		}
	}
}

func readEmbeddedMap(f *bytes.Buffer) {
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

func readSettings(f *bytes.Buffer) {
	b := noerr(wznet.ReadBytes(f, int(noerr(wznet.ReadUBE32(f)))))
	var s ReplaySettings
	must(json.Unmarshal(b, &s))
	if s.ReplayFormatVer != 2 {
		log.Printf("Replay format version is odd: %d, should be 2", s.ReplayFormatVer)
	}
	log.Printf("Replay version %s (netcode %d.%d)", s.GameOptions.VersionString, s.Major, s.Minor)
	log.Printf("Map: %q %q", s.GameOptions.Game.Map, s.GameOptions.Game.Hash)
	mapHash = s.GameOptions.Game.Hash
	if *overrideMapHash != "" {
		mapHash = *overrideMapHash
	}
	replayOptions = s.GameOptions
	log.Printf("Recorded by host: %v", s.GameOptions.NetplayBComms)
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

func readMagic(f *bytes.Buffer) {
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
