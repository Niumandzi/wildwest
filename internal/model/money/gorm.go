package money

type Money struct {
	UserID int `gorm:"primaryKey;unique;column:user_id"`
	Gold   int `gorm:"default:0"`
	Silver int `gorm:"default:0"`
}
