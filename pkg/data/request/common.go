package request

type FieldUpdate struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}
