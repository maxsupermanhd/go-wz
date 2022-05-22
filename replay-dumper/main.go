package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
)

var (
	filepath = flag.String("f", "./replay.wzrp", "Path to replay to dump")
)

func main() {
	log.Println("Replay dumper starting up...")
	flag.Parse()

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
		msgid, ok := netMessageType[pType]
		if !ok {
			log.Printf("Unknown message type: %d", pType)
		} else {
			log.Printf("Message from %2d %-28s len %6d", pPlayer, msgid, len(data))
			// r := bytes.NewReader(data)
			// switch pType {
			// case GAME_RESEARCHSTATUS:
			// 	msg := struct {
			// 		player   uint8
			// 		start    bool
			// 		building uint32
			// 		topic    uint32
			// 	}{}
			// 	must(binary.Read(r, binary.BigEndian, &msg))
			// 	spew.Dump(msg)
			// }
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
	len := noerr(readUBE32(f))
	b := noerr(readBytes(f, int(len)))
	var s ReplaySettings
	must(json.Unmarshal(b, &s))
	if s.ReplayFormatVer != 2 {
		log.Printf("Replay format version is odd: %d, should be 2", s.ReplayFormatVer)
	}
	log.Printf("Replay netcode %d.%d", s.Major, s.Minor)
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
