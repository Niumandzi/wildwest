package gunfight

type QueueRequest struct {
	Gold int `json:"gold" example:"100" extensions:"x-order=2"`
}

type QueueResponse struct {
	ID int `json:"id" example:"1" extensions:"x-order=1"`
}
