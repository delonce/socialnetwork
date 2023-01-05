package service

import "time"

type User struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash string    `json:"-" bson:"password"`
	Email        string    `json:"email" bson:"email"`
	LastEnt      time.Time `json:"time" bson:"time"`
}

type RefreshToken struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UUID      string    `json:"uuid" bson:"uuid"`
	UserID    string    `json:"userID" bson:"userID"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
}

type FriendRequest struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	From       string    `json:"from" bson:"from"`
	To         string    `json:"to" bson:"to"`
	DateAt     time.Time `json:"date" bson:"date"`
	IsAccepted bool      `json:"isaccepted" bson:"isaccepted"`
	IsDenied   bool      `json:"isdenied" bson:"isdenied"`
}

type Message struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	From      string    `json:"from" bson:"from"`
	To        string    `json:"to" bson:"to"`
	Text      string    `json:"text" bson:"text"`
	DateAt    time.Time `json:"date" bson:"date"`
	IsChecked bool      `json:"ischecked" bson:"ischecked"`
}

type ViewMessage struct {
	From       string
	To         string
	Text       string
	FormatDate string
	IsChecked  bool
}

type ViewDialog struct {
	From         string    `json:"from" bson:"from"`
	To           string    `json:"to" bson:"to"`
	Text         string    `json:"text" bson:"text"`
	DateAt       time.Time `json:"date" bson:"date"`
	IsChecked    bool      `json:"ischecked" bson:"ischecked"`
	AmountNewMsg int64
}
