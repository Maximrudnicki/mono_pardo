package request

type CreateSetRequest struct {
	UserId int
	Name   string `json:"name"`
}

type GetSetRequest struct {
	UserId  int
	GroupId string `json:"group_id"`
}

type GetSetsRequest struct {
	UserId int
}

type UpdateSetRequest struct {
	UserId  int
	GroupId string        `json:"group_id"`
	Updates []FieldUpdate `json:"updates"`
}

type DeleteSetRequest struct {
	UserId  int
	GroupId string `json:"group_id"`
}

type AddWordRequest struct {
	UserId  int
	GroupId string
	WordId  int `json:"word_id"`
}

type RemoveWordRequest struct {
	UserId  int
	GroupId string
	WordId  int `json:"word_id"`
}
