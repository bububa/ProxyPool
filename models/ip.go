package models

// IP struct
type IP struct {
	Data string `bson:"_id" json:"ip"`
	Type string `bson:"type" json:"type"`
}
