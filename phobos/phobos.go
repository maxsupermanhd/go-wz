package phobos

import (
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
)

var (
	ErrBadPhobosAnswer = errors.New("phobos returned not 200 status code")
	ErrMultipleByHash  = errors.New("phobos responed with multiple maps")
	ErrNotFound        = errors.New("map is not found on phobos")
)

type MapInfo struct {
	MapName         string `json:"map_name"`
	MapSha          string `json:"map_sha"`
	MapSize         string `json:"map_size"`
	MapPlayers      int    `json:"map_players"`
	MapSpectators   int    `json:"map_spectators"`
	MapScavengers   int    `json:"map_scavengers"`
	MapOilwells     int    `json:"map_oilwells"`
	MapAuthor       string `json:"map_author"`
	MapLicense      string `json:"map_license"`
	MapType         string `json:"map_type"`
	MapScrX1        int    `json:"map_scr_x1"`
	MapScrY1        int    `json:"map_scr_y1"`
	MapScrX2        int    `json:"map_scr_x2"`
	MapScrY2        int    `json:"map_scr_y2"`
	MapAsymmetrical bool   `json:"map_asymmetrical"`
}

func FetchOnePhobosInfo(hash string) (MapInfo, error) {
	var ret MapInfo
	ret.MapSha = hash
	resp, err := http.Get("https://wz2100.euphobos.ru/maps/?api=json&s=" + hash)
	if err != nil {
		return ret, err
	}
	if resp.StatusCode != 200 {
		return ret, ErrBadPhobosAnswer
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	res := map[string]MapInfo{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return ret, err
	}
	if len(res) > 1 {
		return ret, ErrMultipleByHash
	}
	for _, v := range res {
		return v, nil
	}
	return ret, ErrNotFound
}

type PreviewType string

const (
	PreviewTypePixelPerfect = PreviewType(".png")
	PreviewTypeLargeJPEG    = PreviewType(".jpg")
	PreviewTypeHeightmap    = PreviewType("_h.png")
)

func FetchMapPreview(hash string, t PreviewType) (image.Image, error) {
	url := "https://wz2100.euphobos.ru/maps/preview/" + hash + string(t)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, ErrBadPhobosAnswer
	}
	var i image.Image
	switch t {
	case PreviewTypeHeightmap:
		fallthrough
	case PreviewTypePixelPerfect:
		i, err = png.Decode(resp.Body)
	case PreviewTypeLargeJPEG:
		i, err = jpeg.Decode(resp.Body)
	}
	if err != nil {
		return nil, err
	}
	return i, err
}
