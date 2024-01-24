package model

import "github.com/google/uuid"

type Deck struct {
	ID        uuid.UUID `json:"deck_id" gorm:"primaryKey"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     []Card    `json:"cards" gorm:"foreignKey:DeckID"`
}

type Card struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DeckID uuid.UUID `json:"deck_id" gorm:"index;foreignKey:DeckID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Value  string    `json:"value"`
	Suit   string    `json:"suit"`
	Code   string    `json:"code"`
}
