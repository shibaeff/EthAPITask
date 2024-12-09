package models

type BlockReward struct {
	Status bool  `json:"status"`
	Reward int64 `json:"reward"`
}

type SyncDuties struct {
	Validators []string `json:"validators"`
}

type Error struct {
	Error string `json:"error"`
}
