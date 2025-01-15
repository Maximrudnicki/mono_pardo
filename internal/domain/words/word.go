package words

import (
	"errors"
	"strings"
	"time"
)

type Word struct {
	Id         int       `gorm:"type:int;primary_key"`
	Word       string    `gorm:"type:varchar;not null;uniqueIndex:idx_user_word,priority:2,column:word"`
	Definition string    `gorm:"type:varchar;not null"`
	UserId     int       `gorm:"not null;uniqueIndex:idx_user_word,priority:1,column:user_id"`
	CreatedAt  time.Time `gorm:"default:now()"`

	IsLearned       bool `gorm:"default:false"` // status of the word
	Cards           bool `gorm:"default:false"`
	WordTranslation bool `gorm:"default:false"`
	Constructor     bool `gorm:"default:false"`
	WordAudio       bool `gorm:"default:false"`
}

func NewWord(word, definition string, userId int) (*Word, error) {
	if strings.TrimSpace(word) == "" {
		return nil, errors.New("word is required field")
	}
	if strings.TrimSpace(definition) == "" {
		return nil, errors.New("definition is required field")
	}
	if userId <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return &Word{
		Word:       strings.TrimSpace(word),
		Definition: strings.TrimSpace(definition),
		UserId:     userId,
	}, nil
}
