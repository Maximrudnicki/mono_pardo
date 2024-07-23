package request

type AddWordToStudentRequest struct {
	Word       string `json:"word,omitempty"`
	Definition string `json:"definition,omitempty"`
	UserId     int    `json:"user_id,omitempty"`
}

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
	Token      string `json:"token"`
	WordId     int    `json:"word_id"`
	Definition string `json:"definition"`
}

type UpdateWordStatusRequest struct {
	Token     string `json:"token"`
	WordId    int    `json:"word_id"`
	IsLearned bool   `json:"is_learned"`
}

type ManageTrainingsRequest struct {
	Token          string `json:"token"`
	TrainingResult bool   `json:"result"`
	Training       string `json:"training"`
	WordId         int    `json:"word_id"`
}
