package response

import "time"

type SetResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Words     []int     `json:"words"`
}
