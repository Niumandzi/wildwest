package money

// BaseResponse represents the money model
// @Description This is a money model
type BaseResponse struct {
	UserID int `json:"user_id" example:"1" extensions:"x-order=1"`
	Gold   int `json:"gold" example:"100" extensions:"x-order=2"`
	Silver int `json:"silver" example:"50" extensions:"x-order=3"`
}
