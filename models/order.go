package models

type RoomType string

const (
	Econom   RoomType = "econom"
	Standart RoomType = "standart"
	Lux      RoomType = "lux"
)

type Order struct {
	UserID    string   `json:"user_id"`
	UserEmail string   `json:"user_email"`
	RoomType  RoomType `json:"room_type"`
	RoomIDs   []string `json:"room_ids"`
	From      string   `json:"from"`
	To        string   `json:"to"`
}
