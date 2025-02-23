package main

type User struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"type:varchar(100);not null;default:''"`

	// belongs to: many -> one
	CompanyID int
	Company   Company `gorm:"foreignkey:CompanyID"`

	// has one: one -> one
	Card Card `gorm:"foreignkey:UserID"`

	// has one self: one -> one
	ManagerID int
	Manager   *User `gorm:"foreignkey:ManagerID"`

	// has many: one -> many
	Posts []Post `gorm:"foreignkey:UserID"`

	// many to many: many <-> many
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Card struct {
	ID     int    `gorm:"primary_key"`
	Number string `gorm:"type:varchar(100);not null;default:''"`
	UserID int
}

type Post struct {
	ID     int    `gorm:"primary_key"`
	Title  string `gorm:"type:varchar(100);not null;default:''"`
	UserID int
}

type Language struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"type:varchar(100);not null;default:''"`
}
