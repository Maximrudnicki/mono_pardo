package request

type CreateWordRequest struct {
	Token      string `json:"token"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type DeleteWordRequest struct {
	Token  string `json:"token"`
	WordId int    `json:"word_id"`
}

type FindWordRequest struct {
	WordId int `json:"word_id"`
}

type VocabRequest struct {
	TokenType string `json:"token_type"` // Bearer
	Token     string `json:"token"`
}

type UpdateWordRequest struct {
	Token string       `json:"token"`
	Words []WordUpdate `json:"words"`
}

type WordUpdate struct {
	WordId  int           `json:"id"`
	Updates []FieldUpdate `json:"updates"`
}

type FieldUpdate struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type UpdateWordStatusRequest struct {
	Token     string `json:"token"`
	WordId    int    `json:"word_id"`
	IsLearned bool   `json:"is_learned"`
}
