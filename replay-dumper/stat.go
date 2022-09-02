package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"sort"

	"github.com/maxsupermanhd/go-wz/wznet"
)

var (
	aResearch     = []IDNameObject{}
	aBody         = []IDNameObject{}
	aConstruction = []IDNameObject{}
	aECM          = []IDNameObject{}
	aPropulsion   = []IDNameObject{}
	aRepair       = []IDNameObject{}
	aSensor       = []IDNameObject{}
	aStructure    = []IDNameObject{}
	aWeapons      = []IDNameObject{}
)

type IDNameObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func loadDefaultData(p string, out *[]IDNameObject) (err error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	r := map[string]IDNameObject{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ret := []IDNameObject{}
	for _, k := range keys {
		ret = append(ret, r[k])
	}
	*out = ret
	return nil
}

func loadStatsData(p string) (err error) {
	if p == "" {
		p = "./data/mp/stats/"
	}
	if err := loadDefaultData(path.Join(p, "body.json"), &aBody); err != nil {
		return err
	}
	if err := loadDefaultData(path.Join(p, "construction.json"), &aConstruction); err != nil {
		return err
	}
	if err := loadDefaultData(path.Join(p, "ecm.json"), &aECM); err != nil {
		return err
	}
	if err := loadDefaultData(path.Join(p, "propulsion.json"), &aPropulsion); err != nil {
		return err
	}
	if err := loadDefaultData(path.Join(p, "repair.json"), &aRepair); err != nil {
		return err
	}
	if err := loadDefaultData(path.Join(p, "research.json"), &aResearch); err != nil {
		return err
	}
	if err := loadDefaultData(path.Join(p, "sensor.json"), &aSensor); err != nil {
		return err
	}
	if err := loadDefaultData(path.Join(p, "structure.json"), &aStructure); err != nil {
		return err
	}
	if err := loadDefaultData(path.Join(p, "weapons.json"), &aWeapons); err != nil {
		return err
	}
	return nil
}

type DroidDef struct {
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

type DroidTyped struct {
	Name       string
	ID         uint32
	Type       int32
	Body       string
	Brain      string
	Propulsion string
	Repairunit string
	Ecm        string
	Sensor     string
	Construct  string
}

func isOf(v string, c ...string) bool {
	for _, t := range c {
		if t == v {
			return true
		}
	}
	return false
}

func checkDroidIllegal(d DroidTyped) string {
	if isOf(d.Propulsion, "BaBaLegs", "BaBaProp") {
		return "Ba Ba prop"
	}
	if d.Propulsion == "CyborgLegs" && !isOf(d.Body, "CyborgHeavyBody", "CyborgLightBody") {
		return "Cyborg with wrong body"
	}
	return ""
}

func typifyDroid(d DroidDef) *DroidTyped {
	ret := &DroidTyped{
		Name: d.Name,
		ID:   d.ID,
		Type: d.Type,
	}
	if len(aBody) > int(d.Body) {
		ret.Body = aBody[d.Body].ID
	}
	if len(aPropulsion) > int(d.Propulsion) {
		ret.Propulsion = aPropulsion[d.Propulsion].ID
	}
	if len(aRepair) > int(d.Repairunit) {
		ret.Repairunit = aRepair[d.Repairunit].ID
	}
	if len(aECM) > int(d.Ecm) {
		ret.Ecm = aECM[d.Ecm].ID
	}
	if len(aSensor) > int(d.Sensor) {
		ret.Sensor = aSensor[d.Sensor].ID
	}
	return ret
}

func refToStructName(ref uint32) string {
	if ref&wznet.STAT_MASK == wznet.STAT_STRUCTURE {
		structid := int(ref - wznet.STAT_STRUCTURE)
		if len(aStructure) <= structid {
			log.Printf("Structure ref lookup overflow %d, total %d", structid, len(aStructure))
			return "overflow"
		}
		return aStructure[structid].Name
	}
	return "notastructure"
}
