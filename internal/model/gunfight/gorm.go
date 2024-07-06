package gunfight

import "time"

type Queue struct {
	UserID int `gorm:"unique;not null"`
	Gold   int `gorm:"not null"`
}

type Game struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	User1ID   int       `gorm:"not null"`
	User2ID   int       `gorm:"not null"`
	WinnerID  *int      `gorm:"check:winner_id IS NULL OR winner_id = user_1_id OR winner_id = user_2_id"`
	StartDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EndDate   *time.Time
}
