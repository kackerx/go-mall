package main

type Model interface {
	Key() string
	TableName() string
}

type User struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"type:varchar(100);not null;default:''"`
}

func (u User) Key() string {
	return "id"
}

func (u User) TableName() string {
	return "users"
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

func (u Post) Key() string {
	return "id"
}

func (u Post) TableName() string {
	return "posts"
}

type Language struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"type:varchar(100);not null;default:''"`
}

type Company struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"type:varchar(100);not null;default:''"`
}
