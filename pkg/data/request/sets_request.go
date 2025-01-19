package request

type CreateSetRequest struct {
	UserId int
	Name   string `json:"name"`
}

type GetSetRequest struct {
	UserId    int
	WordSetId string
}

type GetSetsRequest struct {
	UserId int
}

type UpdateSetRequest struct {
	UserId    int
	WordSetId string
	Updates   []FieldUpdate `json:"updates"`
}

type DeleteSetRequest struct {
	UserId    int
	WordSetId string
}

type AddWordRequest struct {
	UserId    int
	WordSetId string
	Words     []int `json:"words"`
}

type RemoveWordRequest struct {
	UserId    int
	WordSetId string
	Words     []int `json:"words"`
}
