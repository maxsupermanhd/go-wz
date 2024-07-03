package mapsdatabase

import (
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"time"
)

// https://github.com/Warzone2100/maps-database/blob/main/docs/API.md#example-1
type MapInfo struct {
	Name    string `json:"name"`
	Slots   int    `json:"slots"`
	Tileset string `json:"tileset"`
	Author  string `json:"author"`
	License string `json:"license"`
	Created string `json:"created"`
	Size    struct {
		W int `json:"w"`
		H int `json:"h"`
	} `json:"size"`
	Scavs    int `json:"scavs"`
	OilWells int `json:"oilWells"`
	Player   struct {
		Units struct {
			Eq  bool `json:"eq"`
			Min int  `json:"min"`
			Max int  `json:"max"`
		} `json:"units"`
		Structs struct {
			Eq  bool `json:"eq"`
			Min int  `json:"min"`
			Max int  `json:"max"`
		} `json:"structs"`
		ResourceExtr struct {
			Eq  bool `json:"eq"`
			Min int  `json:"min"`
			Max int  `json:"max"`
		} `json:"resourceExtr"`
		PwrGen struct {
			Eq  bool `json:"eq"`
			Min int  `json:"min"`
			Max int  `json:"max"`
		} `json:"pwrGen"`
		RegFact struct {
			Eq  bool `json:"eq"`
			Min int  `json:"min"`
			Max int  `json:"max"`
		} `json:"regFact"`
		VtolFact struct {
			Eq  bool `json:"eq"`
			Min int  `json:"min"`
			Max int  `json:"max"`
		} `json:"vtolFact"`
		CyborgFact struct {
			Eq  bool `json:"eq"`
			Min int  `json:"min"`
			Max int  `json:"max"`
		} `json:"cyborgFact"`
		ResearchCent struct {
			Eq  bool `json:"eq"`
			Min int  `json:"min"`
			Max int  `json:"max"`
		} `json:"researchCent"`
		DefStruct struct {
			Eq  bool `json:"eq"`
			Min int  `json:"min"`
			Max int  `json:"max"`
		} `json:"defStruct"`
	} `json:"player"`
	Hq       [][2]int `json:"hq"`
	Download struct {
		Type     string `json:"type"`
		Repo     string `json:"repo"`
		Path     string `json:"path"`
		Uploaded string `json:"uploaded"`
		Hash     string `json:"hash"`
		Size     int    `json:"size"`
	} `json:"download"`
}

var defaultClient = &http.Client{
	Timeout: 5 * time.Second,
}

func FetchMapInfo(hash string) (*MapInfo, error) {
	return FetchMapInfoWithClient(hash, defaultClient)
}

func FetchMapInfoWithClient(hash string, cl *http.Client) (*MapInfo, error) {
	if cl == nil {
		cl = defaultClient
	}
	ir, err := cl.Get("https://maps.wz2100.net/api/v1/maps/" + hash + "/info.json")
	if err != nil {
		return nil, err
	}
	var info MapInfo
	err = json.NewDecoder(ir.Body).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func FetchMapBlob(hash string) ([]byte, error) {
	return FetchMapBlobWithClient(hash, defaultClient)
}

func FetchMapBlobWithClient(hash string, cl *http.Client) ([]byte, error) {
	if cl == nil {
		cl = defaultClient
	}
	info, err := FetchMapInfoWithClient(hash, cl)
	if err != nil {
		return nil, err
	}
	br, err := cl.Get("https://github.com/Warzone2100/maps-" + info.Download.Repo + "/releases/download/" + info.Download.Path)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(br.Body)
}

func FetchMapPreview(hash string) (image.Image, error) {
	return FetchMapPreviewWithClient(hash, defaultClient)
}

func FetchMapPreviewWithClient(hash string, cl *http.Client) (image.Image, error) {
	if cl == nil {
		cl = defaultClient
	}
	br, err := cl.Get("https://maps.wz2100.net/api/v1/maps/" + hash + "/preview.png")
	if err != nil {
		return nil, err
	}
	return png.Decode(br.Body)
}
