package response

import "time"

type SetResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Words     []int     `json:"words"`
}
