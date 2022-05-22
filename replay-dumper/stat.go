package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"sort"
)

type sResearch struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func loadResearchData(p string) (ret []sResearch, err error) {
	if p == "" {
		p = "./data/mp/stats/"
	}
	b, err := os.ReadFile(path.Join(p, "research.json"))
	if err != nil {
		return ret, err
	}
	r := map[string]sResearch{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return ret, err
	}
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		ret = append(ret, r[k])
	}
	log.Printf("Loaded %d researches", len(keys))
	return ret, nil
}
