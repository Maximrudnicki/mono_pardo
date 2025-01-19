package request

type CreateSetRequest struct {
	UserId int
	Name   string `json:"name"`
}

type GetSetRequest struct {
	UserId    int
	WordSetId string `json:"word_set_id"`
}

type GetSetsRequest struct {
	UserId int
}

type UpdateSetRequest struct {
	UserId    int
	WordSetId string        `json:"word_set_id"`
	Updates   []FieldUpdate `json:"updates"`
}

type DeleteSetRequest struct {
	UserId    int
	WordSetId string `json:"word_set_id"`
}

type AddWordRequest struct {
	UserId    int
	WordSetId string
	WordId    int `json:"word_id"`
}

type RemoveWordRequest struct {
	UserId    int
	WordSetId string
	WordId    int `json:"word_id"`
}
