package model

import "time"

type Word struct {
	Id         int       `gorm:"type:int;primary_key"`
	Word       string    `gorm:"type:varchar;not null"`
	Definition string    `gorm:"type:varchar;not null"`
	UserId     int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"default:now()"`

	IsLearned       bool `gorm:"default:false"` // status of the word
	Cards           bool `gorm:"default:false"`
	WordTranslation bool `gorm:"default:false"`
	Constructor     bool `gorm:"default:false"`
	WordAudio       bool `gorm:"default:false"`
}
