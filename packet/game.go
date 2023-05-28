package packet

import (
	"fmt"
	"io"
	"path"
	"runtime"

	"github.com/maxsupermanhd/go-wz/wznet"
)

type NetPacket interface {
	Name() string
	Type() byte
	Length() uint32
}

// ParsePacket will panic if fails to parse packet, will error when unknown packet type
func ParsePacket(pt byte, l uint32, r io.Reader) (ret NetPacket, err error) {
	p := packetParsers[pt]
	if p == nil {
		pkname, ok := wznet.NetMessageType[pt]
		if ok {
			return nil, fmt.Errorf("no packet parser for %v", pkname)
		} else {
			return nil, fmt.Errorf("no packet parser for unknown packet type %v", pt)
		}
	}
	ret = p(pk{pt, l}, r)
	return ret, err
}

type packetParser func(pk, io.Reader) NetPacket

var (
	packetParsers = map[byte]packetParser{
		wznet.GAME_GAME_TIME:      ParseGameGameTime,
		wznet.GAME_STRUCTUREINFO:  ParseGameStructInfo,
		wznet.GAME_RESEARCHSTATUS: ParseGameResearchStatus,
		wznet.GAME_DROIDINFO:      ParseGameDroidInfo,
		wznet.GAME_PLAYER_LEFT:    ParseGamePlayerLeft,
		wznet.GAME_GIFT:           ParseGameGift,
		wznet.GAME_LASSAT:         ParseGameLasSat,
		wznet.REPLAY_ENDED:        ParseNothing,
	}
)

type pk struct {
	t byte
	l uint32
}

func (p pk) Name() string {
	return wznet.NetMessageType[p.t]
}
func (p pk) Type() byte {
	return p.t
}
func (p pk) Length() uint32 {
	return p.l
}

func ParseNothing(p pk, r io.Reader) NetPacket {
	return p
}

type PkGameGameTime struct {
	pk
	LatencyTicks  uint32
	GameTime      uint32
	CRC           uint16
	WantedLatency uint16
}

func ParseGameGameTime(p pk, r io.Reader) NetPacket {
	ret := PkGameGameTime{pk: p}
	ret.LatencyTicks = panicerr(wznet.NETreadU32(r))
	ret.GameTime = panicerr(wznet.NETreadU32(r))
	ret.CRC = panicerr(wznet.NETreadU16(r))
	ret.WantedLatency = panicerr(wznet.NETreadU16(r))
	return ret
}

type PkGameStructInfo struct {
	pk
	Player       uint8
	StructID     uint32
	StructInfo   uint8
	Droid        PkGameStructInfoDroidDef
	DroidWeapons []uint32
}
type PkGameStructInfoDroidDef struct {
	Name       string
	ID         uint32
	Type       int32
	Body       uint8
	Brain      uint8
	Propulsion uint8
	Repairunit uint8
	Ecm        uint8
	Sensor     uint8
	Construct  uint8
}

func ParseGameStructInfo(p pk, r io.Reader) NetPacket {
	ret := PkGameStructInfo{pk: p}
	ret.Player = panicerr(wznet.NETreadU8(r))
	ret.StructID = panicerr(wznet.NETreadU32(r))
	ret.StructInfo = panicerr(wznet.NETreadU8(r))
	if ret.StructInfo == wznet.STRUCTUREINFO_MANUFACTURE {
		ret.Droid = PkGameStructInfoDroidDef{
			panicerr(wznet.NETstring(r)),  // droid Name
			panicerr(wznet.NETreadU32(r)), // droid ID
			panicerr(wznet.NETreadS32(r)), // droid Type
			panicerr(wznet.NETreadU8(r)),  // droid Body
			panicerr(wznet.NETreadU8(r)),  // droid Brain
			panicerr(wznet.NETreadU8(r)),  // droid Propulsion
			panicerr(wznet.NETreadU8(r)),  // droid Repairunit
			panicerr(wznet.NETreadU8(r)),  // droid Ecm
			panicerr(wznet.NETreadU8(r)),  // droid Sensor
			panicerr(wznet.NETreadU8(r)),  // droid Construct
		}
		droidNumWeapons := panicerr(wznet.NETreadU8(r))
		for i := uint8(0); i < droidNumWeapons; i++ {
			ret.DroidWeapons = append(ret.DroidWeapons, panicerr(wznet.NETreadU32(r)))
		}
	}
	return ret
}

type PkGameResearchStatus struct {
	pk
	Player   uint8
	Start    bool
	Building uint32
	Topic    uint32
}

func ParseGameResearchStatus(p pk, r io.Reader) NetPacket {
	ret := PkGameResearchStatus{pk: p}
	ret.Player = panicerr(wznet.NETreadU8(r))
	if panicerr(wznet.NETreadU8(r)) > 0 {
		ret.Start = true
	}
	ret.Building = panicerr(wznet.NETreadU32(r))
	ret.Topic = panicerr(wznet.NETreadU32(r))
	return ret
}

type PkGameDroidInfo struct {
	pk
	Player    uint8
	SubType   wznet.DroidOrderSybType
	Order     wznet.DORDER
	DestID    uint32
	DestType  uint32
	CoordX    int32
	CoordY    int32
	StructRef uint32
	Direction uint16
	CoordX2   int32
	CoordY2   int32
	Index     uint32
	Add       bool
	SecOrder  uint32
	SecState  uint32
	Droids    []uint32
}

func ParseGameDroidInfo(p pk, r io.Reader) NetPacket {
	ret := PkGameDroidInfo{pk: p}
	ret.Player = panicerr(wznet.NETreadU8(r))
	ret.SubType = wznet.DroidOrderSybType(panicerr(wznet.NETreadU32(r)))
	switch ret.SubType {
	case wznet.DroidOrderSybTypeObj:
		fallthrough
	case wznet.DroidOrderSybTypeLoc:
		ret.Order = wznet.DORDER(panicerr(wznet.NETreadU32(r)))
		if ret.SubType == wznet.DroidOrderSybTypeObj {
			ret.DestID = panicerr(wznet.NETreadU32(r))
			ret.DestType = panicerr(wznet.NETreadU32(r))
		} else {
			ret.CoordX = panicerr(wznet.NETreadS32(r))
			ret.CoordY = panicerr(wznet.NETreadS32(r))
		}
		if ret.Order == wznet.DORDER_BUILD || ret.Order == wznet.DORDER_LINEBUILD {
			ret.StructRef = panicerr(wznet.NETreadU32(r))
			ret.Direction = panicerr(wznet.NETreadU16(r))
		}
		if ret.Order == wznet.DORDER_LINEBUILD {
			ret.CoordX2 = panicerr(wznet.NETreadS32(r))
			ret.CoordY2 = panicerr(wznet.NETreadS32(r))
		}
		if ret.Order == wznet.DORDER_BUILDMODULE {
			ret.Index = panicerr(wznet.NETreadU32(r))
		}
		if panicerr(wznet.NETreadU8(r)) > 0 {
			ret.Add = true
		}
	case wznet.DroidOrderSybTypeSec:
		ret.SecOrder = panicerr(wznet.NETreadU32(r))
		ret.SecState = panicerr(wznet.NETreadU32(r))
	}
	num := panicerr(wznet.NETreadU32(r))
	droiddelta := uint32(0)
	for i := uint32(0); i < num; i++ {
		droiddelta += panicerr(wznet.NETreadU32(r))
		ret.Droids = append(ret.Droids, droiddelta)
	}
	return ret
}

type PkGamePlayerLeft struct {
	pk
	Player uint8
}

func ParseGamePlayerLeft(p pk, r io.Reader) NetPacket {
	ret := PkGamePlayerLeft{pk: p}
	ret.Player = panicerr(wznet.NETreadU8(r))
	return ret
}

type PkGameGift struct {
	pk
	GiftType wznet.GIFT_TYPE
	From     uint8
	To       uint8
	DroidID  uint32
}

func ParseGameGift(p pk, r io.Reader) NetPacket {
	ret := PkGameGift{pk: p}
	ret.GiftType = wznet.GIFT_TYPE(panicerr(wznet.NETreadU8(r)))
	ret.From = panicerr(wznet.NETreadU8(r))
	ret.To = panicerr(wznet.NETreadU8(r))
	ret.DroidID = panicerr(wznet.NETreadU32(r))
	return ret
}

type PkGameLasSat struct {
	pk
	Player       uint8
	ID           uint32
	TargetID     uint32
	TargetPlayer uint8
}

func ParseGameLasSat(p pk, r io.Reader) NetPacket {
	ret := PkGameLasSat{pk: p}
	ret.Player = panicerr(wznet.NETreadU8(r))
	ret.ID = panicerr(wznet.NETreadU32(r))
	ret.TargetID = panicerr(wznet.NETreadU32(r))
	ret.TargetPlayer = panicerr(wznet.NETreadU8(r))
	return ret
}

func panicerr[T any](t T, err error) T {
	if err != nil {
		pc, filename, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("Error: %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), path.Base(filename), line, err))
	}
	return t
}
