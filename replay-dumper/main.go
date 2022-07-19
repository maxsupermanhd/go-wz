package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/maxsupermanhd/go-wz/wznet"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

var (
	filepath       = flag.String("f", "./replay.wzrp", "Path to replay to dump")
	statsdir       = flag.String("stats", "./data/mp/stats/", "Path to stats directory")
	mapout         = flag.String("mapout", "./map.wz", "Path to save embedded map. Use - to disable")
	aResearch      = []sResearch{}
	netPlayPlayers = []NetplayPlayers{}
	namePadLength  = 2
)

func main() {
	log.SetFlags(0)
	log.Println("Replay dumper starting up...")
	flag.Parse()

	log.Printf("Loading research from [%s]...", *statsdir)
	aResearch = noerr(loadResearchData(*statsdir))

	log.Printf("Opening file [%s]...", *filepath)
	f := noerr(os.Open(*filepath))

	log.Println("Reading magic...")
	readMagic(f)

	log.Println("Reading header JSON...")
	readSettings(f)

	log.Println("Reading embedded map data...")
	readEmbeddedMap(f)

	log.Printf("Begin of netmessages...")
	readNetMessages(f)

	log.Println("Bye!")
	f.Close()
}

func readNetMessages(f *os.File) {
	total := 0
	gameTime := uint32(0)
	tdata := make([]opts.LineData, 0)
	ldata := make([]opts.LineData, 0)
	labels := []string{}
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line example in Westeros theme",
			Subtitle: "Line chart rendered by the http server this time",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{Type: "inside"}))

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
				latencyTicks := noerr(wznet.NETreadU32(r))
				gameTime = noerr(wznet.NETreadU32(r))
				_ = noerr(wznet.NETreadU16(r))
				wantedLatency := noerr(wznet.NETreadU16(r))
				tdata = append(tdata, opts.LineData{Value: latencyTicks})
				ldata = append(ldata, opts.LineData{Value: wantedLatency})
				labels = append(labels, fmt.Sprint(gameTime))
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
				player := noerr(wznet.NETreadU8(r))
				structID := noerr(wznet.NETreadU32(r))
				structInfo := noerr(wznet.NETreadU8(r))
				if structInfo == 0 {
					name := noerr(wznet.NETstring(r))
					droidID := noerr(wznet.NETreadU32(r))
					droidType := noerr(wznet.NETreadS32(r))
					droidBody := noerr(wznet.NETreadU8(r))
					droidBrain := noerr(wznet.NETreadU8(r))
					droidPropulsion := noerr(wznet.NETreadU8(r))
					droidRepairunit := noerr(wznet.NETreadU8(r))
					droidEcm := noerr(wznet.NETreadU8(r))
					droidSensor := noerr(wznet.NETreadU8(r))
					droidConstruct := noerr(wznet.NETreadU8(r))
					droidNumWeapons := noerr(wznet.NETreadU8(r))
					droidWeapons := []uint32{}
					for i := uint8(0); i < droidNumWeapons; i++ {
						droidWeapons = append(droidWeapons, noerr(wznet.NETreadU32(r)))
					}
					rprint(`GAME_STRUCTUREINFO player %v struct %v info %v name [%v]
id %v
type %v
body %v
brain %v
prop %v
rep %v
ecm %v
sensor %v
construct %v
weapons %v`, player, structID, structInfo, name, droidID,
						droidType,
						droidBody,
						droidBrain,
						droidPropulsion,
						droidRepairunit,
						droidEcm,
						droidSensor,
						droidConstruct, droidWeapons)
				}
			case wznet.GAME_DROIDINFO:
				// player := noerr(wznet.NETreadU8(r))
				// if player != pPlayer {
				// 	log.Printf("Player missmatch in %s (%d netmessage %d packet)", msgid, pPlayer, player)
				// }
				// subtype := noerr(wznet.NETreadU32(r))
				// switch subtype {
				// case DroidOrderSybTypeObj:
				// 	fallthrough
				// case DroidOrderSybTypeLoc:
				// 	order := noerr(wznet.NETreadU32(r))
				// 	if subtype == DroidOrderSybTypeObj {
				// 		_ = noerr(wznet.NETreadU32(r)) // destID
				// 		_ = noerr(wznet.NETreadU32(r)) // destType
				// 	} else {
				// 		_ = noerr(wznet.NETreadS32(r))
				// 		_ = noerr(wznet.NETreadS32(r))
				// 	}
				// 	if order == DORDER_BUILD || order == DORDER_LINEBUILD {
				// 		_ = noerr(wznet.NETreadU32(r)) // structref
				// 		_ = noerr(wznet.NETreadU16(r)) // direction
				// 	}
				// 	if order == DORDER_LINEBUILD {
				// 		_ = noerr(wznet.NETreadS32(r)) // pos2 x
				// 		_ = noerr(wznet.NETreadS32(r)) // pos2 y
				// 	}
				// 	if order == DORDER_BUILDMODULE {
				// 		_ = noerr(wznet.NETreadU32(r)) // index
				// 	}
				// 	_ = noerr(wznet.NETreadU8(r)) // add
				// case DroidOrderSybTypeSec:
				// 	_ = noerr(wznet.NETreadU32(r)) // sec order
				// 	_ = noerr(wznet.NETreadU32(r)) // sec state
				// }
			case wznet.GAME_RESEARCHSTATUS:
				player := noerr(wznet.NETreadU8(r))
				start := noerr(wznet.NETreadU8(r))
				building := noerr(wznet.NETreadU32(r))
				topic := noerr(wznet.NETreadU32(r))
				topicname := fmt.Sprint(topic)
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
				rprint("%s on building %d topic %s", action, building, topicname)
				if r.Len() != 0 {
					spew.Dump(data)
					log.Printf("Did not parsed %d bytes", r.Len())
				}
			}
		}
		total++
	}
	var buf bytes.Buffer
	line.SetXAxis(labels).
		// AddSeries("wantedLatency", ldata).
		AddSeries("latencyTicks", tdata).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(&buf)
	must(os.WriteFile("chart.html", buf.Bytes(), 0644))
}

func readEmbeddedMap(f *os.File) {
	dv := noerr(wznet.ReadUBE32(f))
	if dv != 1 {
		log.Printf("Embedded map data version is odd: %d, should be 1", dv)
	}
	len := noerr(wznet.ReadUBE32(f))
	if len != 0 {
		log.Printf("Embedded map data len %d", len)
		b := noerr(wznet.ReadBytes(f, int(len)))
		if *mapout != "-" {
			must(os.WriteFile(*mapout, b, 0644))
		}
	} else {
		log.Println("Embedded map data is empty")
	}
}

func readSettings(f *os.File) {
	b := noerr(wznet.ReadBytes(f, int(noerr(wznet.ReadUBE32(f)))))
	var s ReplaySettings
	must(json.Unmarshal(b, &s))
	if s.ReplayFormatVer != 2 {
		log.Printf("Replay format version is odd: %d, should be 2", s.ReplayFormatVer)
	}
	log.Printf("Replay netcode %d.%d", s.Major, s.Minor)
	log.Printf("Replay version %s", s.GameOptions.VersionString)
	log.Println("Players:")
	for i, p := range s.GameOptions.NetplayPlayers {
		if p.Allocated {
			log.Printf("pos %2d inx %2d name [%s]", p.Position, i, p.Name)
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
	log.Printf("Magic is valid")
}

func gameTimeToString(gt uint32) string {
	return (time.Duration(int(gt/1000)) * time.Second).String()
}
