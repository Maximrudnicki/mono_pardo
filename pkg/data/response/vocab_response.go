package response

import "time"

type VocabResponse struct {
	Id              int       `json:"id"`
	Word            string    `json:"word"`
	Definition      string    `json:"definition"`
	CreatedAt       time.Time `json:"created_at"`
	IsLearned       bool      `json:"is_learned"`
	Cards           bool      `json:"cards"`
	WordTranslation bool      `json:"word_translation"`
	Constructor     bool      `json:"constructor"`
	WordAudio       bool      `json:"word_audio"`
}
