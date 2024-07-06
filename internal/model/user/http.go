package user

type BaseRequest struct {
	ID        int    `json:"id" example:"1" extensions:"x-order=1"`
	FirstName string `json:"first_name" extensions:"x-order=2"`
	LastName  string `json:"last_name" extensions:"x-order=3"`
	Hash      string `json:"hash" extensions:"x-order=4"`
}

// BaseResponse represents the user model
// @Description This is a user model
type BaseResponse struct {
	ID        int    `json:"id" example:"1" extensions:"x-order=1"`
	FirstName string `json:"first_name" extensions:"x-order=2"`
	LastName  string `json:"last_name" extensions:"x-order=3"`
	Link      string `json:"link" example:"dsfb434kdfp" extensions:"x-order=5"`
}
