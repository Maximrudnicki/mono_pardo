package request

type CreateWordRequest struct {
	UserId     int
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type DeleteWordRequest struct {
	UserId int
	WordId int `json:"word_id"`
}

type FindWordRequest struct {
	WordId int `json:"word_id"`
}

type VocabRequest struct {
	UserId int
}

type UpdateWordRequest struct {
	UserId int
	Words  []WordUpdate `json:"words"`
}

type WordUpdate struct {
	WordId  int           `json:"id"`
	Updates []FieldUpdate `json:"updates"`
}

type FieldUpdate struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}
