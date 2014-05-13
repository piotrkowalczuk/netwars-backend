package main

type Stream struct {
	UserId int64 `db:"user_id" json:"id"`
	Type int64 `db:"streamtype" json:"type"`
	Identifier string `db:"handle" json:"identifier"`
	WingsOfLiberty string `db:"wol" json:"wingsOfLiberty"`
	HeartOfTheSwarm string `db:"hots" json:"heartOfTheSwarm"`
	LeagueOfLegends string `db:"lol" json:"leagueOfLegends"`
	BroodWar string `db:"bw" json:"broodWar"`
	Other string `db:"other" json:"other"`
}

type StreamRequest struct {
	Stream
}

func (cusr *StreamRequest) isValid() bool {
	return true
}
