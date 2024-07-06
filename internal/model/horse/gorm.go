package horse

type Horse struct {
	UserID   int `gorm:"primaryKey;unique;column:user_id"`
	Level    int `gorm:"default:1"`
	Distance int `gorm:"default:0"`
}
