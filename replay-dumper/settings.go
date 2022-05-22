package main

type ReplaySettings struct {
	GameOptions     GameOptions `json:"gameOptions"`
	Major           int         `json:"major"`
	Minor           int         `json:"minor"`
	ReplayFormatVer int         `json:"replayFormatVer"`
}
type Game struct {
	Alliance          int           `json:"alliance"`
	Base              int           `json:"base"`
	Hash              string        `json:"hash"`
	InactivityMinutes int           `json:"inactivityMinutes"`
	IsMapMod          bool          `json:"isMapMod"`
	IsRandom          bool          `json:"isRandom"`
	Map               string        `json:"map"`
	MapHasScavengers  bool          `json:"mapHasScavengers"`
	MaxPlayers        int           `json:"maxPlayers"`
	ModHashes         []interface{} `json:"modHashes"`
	Name              string        `json:"name"`
	Power             int           `json:"power"`
	Scavengers        int           `json:"scavengers"`
	TechLevel         int           `json:"techLevel"`
	Type              int           `json:"type"`
}
type StructureLimits struct {
	ID    int `json:"id"`
	Limit int `json:"limit"`
}
type Ingame struct {
	Flags           int               `json:"flags"`
	Side            bool              `json:"side"`
	StructureLimits []StructureLimits `json:"structureLimits"`
}
type Multistats struct {
	Identity        string `json:"identity"`
	Losses          int    `json:"losses"`
	Played          int    `json:"played"`
	RecentKills     int    `json:"recentKills"`
	RecentPowerLost int    `json:"recentPowerLost"`
	RecentScore     int    `json:"recentScore"`
	TotalKills      int    `json:"totalKills"`
	TotalScore      int    `json:"totalScore"`
	Wins            int    `json:"wins"`
}
type NetplayPlayers struct {
	Ai              int    `json:"ai"`
	Allocated       bool   `json:"allocated"`
	Colour          int    `json:"colour"`
	Difficulty      int    `json:"difficulty"`
	Faction         int    `json:"faction"`
	Heartattacktime int    `json:"heartattacktime"`
	Heartbeat       bool   `json:"heartbeat"`
	IsSpectator     bool   `json:"isSpectator"`
	Kick            bool   `json:"kick"`
	Name            string `json:"name"`
	Position        int    `json:"position"`
	Ready           bool   `json:"ready"`
	Team            int    `json:"team"`
}
type GameOptions struct {
	DataHash          []interface{}    `json:"dataHash"`
	Game              Game             `json:"game"`
	Ingame            Ingame           `json:"ingame"`
	Multistats        []Multistats     `json:"multistats"`
	NetplayBComms     bool             `json:"netplay.bComms"`
	NetplayHostPlayer int              `json:"netplay.hostPlayer"`
	NetplayPlayers    []NetplayPlayers `json:"netplay.players"`
	RandSeed          int              `json:"randSeed"`
	SelectedPlayer    int              `json:"selectedPlayer"`
	VersionString     string           `json:"versionString"`
}
