package replay

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/maxsupermanhd/go-wz/packet"
	"github.com/maxsupermanhd/go-wz/wznet"
)

var (
	ErrWrongMagic                    = errors.New("wrong magic")
	ErrWrongReplaySettingsVer        = errors.New("wrong replay settings version")
	ErrWrongReplayEmbeddedMapVersion = errors.New("wrong embedded map version")
)

type Replay struct {
	Settings    ReplaySettings
	EmbeddedMap []byte
	Messages    []ReplayPacket
	End         EndChunk
}

func ReadReplay(r io.Reader) (o *Replay, err error) {
	o = &Replay{}

	_, err = readMagic(r)
	if err != nil {
		return nil, err
	}

	o.Settings, err = readSettings(r)
	if err != nil {
		return nil, err
	}

	o.EmbeddedMap, err = readEmbeddedMap(r)
	if err != nil {
		return nil, err
	}

	// Hack! readNetMessage will panic if packet.ParsePacket panics so we catch it here
	defer func() {
		if pan := recover(); pan != nil {
			err = errors.New(pan.(string))
		}
	}()

	for {
		msg, err := readNetMessage(r)
		if err != nil {
			return nil, err
		}
		o.Messages = append(o.Messages, *msg)
		if msg.Type() == wznet.REPLAY_ENDED {
			break
		}
	}

	o.End, err = readEndChunk(r)
	if err != nil {
		return nil, err
	}

	return o, err
}

func readEndChunk(r io.Reader) (EndChunk, error) {
	var ret EndChunk
	l, err := wznet.ReadUBE32(r)
	if err != nil {
		return ret, err
	}
	b := make([]byte, l)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(b, &ret)
	wznet.ReadUBE32(r)
	return ret, err
}

type ReplayPacket struct {
	Player byte
	packet.NetPacket
}

func readNetMessage(r io.Reader) (*ReplayPacket, error) {
	ret := &ReplayPacket{}
	h := make([]byte, 2)
	_, err := io.ReadFull(r, h)
	if err != nil {
		return nil, err
	}
	ret.Player = h[0]

	log.Println(h)

	l, err := wznet.NETreadU32(r)
	if err != nil {
		return nil, err
	}
	ret.NetPacket, err = packet.ParsePacket(h[1], l, io.LimitReader(r, int64(l)))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func readEmbeddedMap(r io.Reader) ([]byte, error) {
	dv, err := wznet.ReadUBE32(r)
	if err != nil {
		return nil, err
	}
	if dv != 1 {
		return nil, ErrWrongReplayEmbeddedMapVersion
	}
	ml, err := wznet.ReadUBE32(r)
	if err != nil {
		return nil, err
	}
	var b []byte
	if ml != 0 {
		b, err = wznet.ReadBytes(r, int(ml))
	}
	return b, err
}

func readSettings(r io.Reader) (ReplaySettings, error) {
	var s ReplaySettings
	sl, err := wznet.ReadUBE32(r)
	if err != nil {
		return s, err
	}
	sb, err := wznet.ReadBytes(r, int(sl))
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(sb, &s)
	if err != nil {
		return s, err
	}
	if s.ReplayFormatVer != 2 {
		return s, ErrWrongReplaySettingsVer
	}
	return s, nil
}

func readMagic(r io.Reader) ([]byte, error) {
	magic := make([]byte, 4)
	_, err := io.ReadAtLeast(r, magic, 4)
	if err != nil {
		return magic, err
	}
	if !bytes.Equal(magic, []byte{'W', 'Z', 'r', 'p'}) {
		return magic, ErrWrongMagic
	}
	return magic, nil
}
