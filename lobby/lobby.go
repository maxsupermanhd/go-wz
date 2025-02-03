package lobby

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type LobbyRoom struct {
	StructVersion  uint32
	GameName       [64]byte
	DW             [2]uint32
	HostIP         [40]byte
	MaxPlayers     uint32
	CurrentPlayers uint32
	DWFlags        [4]uint32
	SecHost        [2][40]byte
	Extra          [157]byte
	Port           uint16
	MapName        [40]byte
	HostName       [40]byte
	Version        [64]byte
	Mods           [255]byte
	VersionMajor   uint32
	VersionMinor   uint32
	Private        uint32
	Pure           uint32
	ModsCount      uint32
	GameID         uint32
	Limits         uint32
	Future1        uint32
	Future2        uint32
}

type LobbyResponse struct {
	MOTD  string
	Rooms []LobbyRoom
	Code  uint32
	Flags uint32
}

const LobbyAddress = "lobby.wz2100.net:9990"

func LobbyLookup() (LobbyResponse, error) {
	return LobbyLookupAddr(LobbyAddress)
}

func LobbyLookupAddr(addr string) (rsp LobbyResponse, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "list\n")

	var count uint32
	err = binary.Read(conn, binary.BigEndian, &count)
	if err != nil {
		return
	}
	for i := uint32(0); i < count; i++ {
		var room LobbyRoom
		err = binary.Read(conn, binary.BigEndian, &room)
		if err != nil {
			return
		}
		rsp.Rooms = append(rsp.Rooms, room)
	}

	err = binary.Read(conn, binary.BigEndian, &rsp.Code)
	if err != nil {
		return
	}

	var motdlen uint32
	err = binary.Read(conn, binary.BigEndian, &motdlen)
	if err != nil {
		return
	}
	motd := make([]byte, motdlen)
	err = binary.Read(conn, binary.BigEndian, &motd)
	if err != nil {
		return
	}
	rsp.MOTD = string(motd)

	err = binary.Read(conn, binary.BigEndian, &rsp.Flags)
	if errors.Is(err, io.EOF) {
		err = nil
		return
	}

	if (rsp.Flags & 1) == 1 {
		rsp.Rooms = []LobbyRoom{}
	}
	err = binary.Read(conn, binary.BigEndian, &count)
	if err != nil {
		return
	}
	for i := uint32(0); i < count; i++ {
		var room LobbyRoom
		err = binary.Read(conn, binary.BigEndian, &room)
		if err != nil {
			return
		}
		rsp.Rooms = append(rsp.Rooms, room)
	}

	return
}
