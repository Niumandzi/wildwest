package user

type User struct {
	ID        int    `gorm:"primaryKey"`
	FirstName string `gorm:"size:255"`
	LastName  string `gorm:"size:255"`
	Link      string `gorm:"size:255;unique"`
}
