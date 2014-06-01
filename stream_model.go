package main

const (
	STREAM_LIFE_TIME int = 300
)

type Stream struct {
	UserId          int64  `db:"user_id" json:"id"`
	Type            int64  `db:"streamtype" json:"type" binding:"required"`
	Identifier      string `db:"handle" json:"identifier" binding:"required"`
	WingsOfLiberty  string `db:"wol" json:"wingsOfLiberty"`
	HeartOfTheSwarm string `db:"hots" json:"heartOfTheSwarm"`
	LeagueOfLegends string `db:"lol" json:"leagueOfLegends"`
	BroodWar        string `db:"bw" json:"broodWar"`
	Other           string `db:"other" json:"other"`
	IsOnline        bool   `json:"isOnline"`
	Current         struct {
		Broadcaster string        `json:"broadcaster"`
		Preview     StructPreview `json:"preview"`
		Id          int64         `json:"id"`
		Viewers     int64         `json:"viewers"`
		Name        string        `json:"name"`
		Game        string        `json:"game"`
		DisplayName string        `json:"displayName"`
	} `json:"current"`
}

type StructPreview struct {
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Large    string `json:"large"`
	Template string `json:"template"`
}
type StreamRequest struct {
	Stream
}

func (cusr *StreamRequest) isValid() bool {
	return true
}

type StreamTwitch struct {
	Links struct {
		Channel string `json:"channel"`
		Self    string `json:"self"`
	} `json:"_links"`
	Stream StreamTwitchDetails `json:"stream"`
}

type StreamTwitchDetails struct {
	Links struct {
		Self string `json:"self"`
	} `json:"_links"`
	Broadcaster string        `json:"broadcaster"`
	Preview     StructPreview `json:"preview"`
	Id          int64         `json:"_id"`
	Viewers     int64         `json:"viewers"`
	Channel     struct {
		DisplayName string `json:"display_name"`
		Links       struct {
			StreamKey     string `json:"stream_key"`
			Editors       string `json:"editors"`
			Subscriptions string `json:"subscriptions"`
			Commercial    string `json:"commercial"`
			Videos        string `json:"videos"`
			Follows       string `json:"follows"`
			Self          string `json:"self"`
			Chat          string `json:"chat"`
			Features      string `json:"features"`
		} `json:"_links"`
		Teams       []interface{} `json:"teams"`
		Status      string        `json:"status"`
		CreatedAt   string        `json:"created_at"`
		Logo        string        `json:"logo"`
		UpdatedAt   string        `json:"updated_at"`
		Mature      interface{}   `json:"mature"`
		VideoBanner interface{}   `json:"video_banner"`
		Id          int64         `json:"_id"`
		Background  string        `json:"background"`
		Banner      string        `json:"banner"`
		Name        string        `json:"name"`
		Url         string        `json:"url"`
		Game        string        `json:"game"`
	} `json:"channel"`
	Name string `json:"name"`
	Game string `json:"game"`
}
