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
)

var (
	filepath       = flag.String("f", "./replay.wzrp", "Path to replay to dump")
	statsdir       = flag.String("stats", "./data/mp/stats/", "Path to stats directory")
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
	for {
		pPlayer, _ := noerr2(readByte(f))
		pType, _ := noerr2(readByte(f))
		l := uint32(0)
		success := false
		b := byte(0)
		n := uint(0)
		for {
			b, success = noerr2(readByte(f))
			r := false
			r, l = decode_uint32_t(b, l, n)
			n++
			if r {
				break
			}
		}
		if !success {
			log.Printf("Not success, len %d", l)
			break
		}
		if pType == REPLAY_ENDED {
			log.Printf("End of net messages. (%d total)", total)
			break
		}
		data := noerr(readBytes(f, int(l)))
		// msgid, ok := netMessageType[pType]
		_, ok := netMessageType[pType]
		if !ok {
			log.Printf("Unknown message type: %d", pType)
		} else {
			r := bytes.NewReader(data)
			// log.Printf("Message from %2d %-28s len %3d %3d", pPlayer, msgid, len(data), l)
			switch pType {
			case GAME_GAME_TIME:
				_ = noerr(NETreadU32(r))
				gameTime = noerr(NETreadU32(r))
				_ = noerr(NETreadU16(r))
				_ = noerr(NETreadU16(r))
			case GAME_RESEARCHSTATUS:
				player := noerr(NETreadU8(r))
				start := noerr(NETreadU8(r))
				building := noerr(NETreadU32(r))
				topic := noerr(NETreadU32(r))
				topicname := fmt.Sprint(topic)
				if topic > 0 && int(topic) < len(aResearch) {
					topicname = aResearch[topic].Name
				} else {
					log.Printf("Topic overflow or underflow, topic %d total %d", topic, len(aResearch))
				}
				if player != pPlayer {
					log.Printf("Player missmatch in GAME_RESEARCHSTATUS (%d netmessage %d packet)", pPlayer, player)
				}
				action := "picked "
				if start == 0 {
					action = "dropped"
				}
				log.Printf("(%7d % 9s) [% *s] %s on building %d topic %s",
					gameTime, gameTimeToString(gameTime),
					-namePadLength, netPlayPlayers[player].Name,
					action, building, topicname)
				if r.Len() != 0 {
					spew.Dump(data)
					log.Printf("Did not parsed %d bytes", r.Len())
				}
			}
		}
		total++
	}
}

func readEmbeddedMap(f *os.File) {
	dv := noerr(readUBE32(f))
	if dv != 1 {
		log.Printf("Embedded map data version is odd: %d, should be 1", dv)
	}
	len := noerr(readUBE32(f))
	if len != 0 {
		log.Printf("Embedded map data len %d", len)
		_ = noerr(readBytes(f, int(len)))
	} else {
		log.Println("Embedded map data is empty")
	}
}

func readSettings(f *os.File) {
	b := noerr(readBytes(f, int(noerr(readUBE32(f)))))
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
