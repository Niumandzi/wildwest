package horse

// BaseResponse represents the horse model
// @Description This is a horse model
type BaseResponse struct {
	UserID   int `json:"user_id" example:"1" extensions:"x-order=1"`
	Level    int `json:"level" example:"100" extensions:"x-order=2"`
	Distance int `json:"distance" example:"50" extensions:"x-order=3"`
	Speed    int `json:"speed" example:"2000" extensions:"x-order=4"`
}

type GameRequest struct {
	Distance int `json:"distance" extensions:"x-order=1"`
}

type GameResponse struct {
	Earned   int  `json:"earned"  extensions:"x-order=1"`
	Record   bool `json:"record" extensions:"x-order=2"`
	Distance int  `json:"distance" extensions:"x-order=3"`
}
