package models

type User struct {
	Name     string `gorm:"size:30;nol null" json:"name"`
	Email    string `gorm:"primaryKey" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Gender   string `gorm:"not null" json:"gender"`
}
