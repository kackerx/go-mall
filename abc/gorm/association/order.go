package main

type Company struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"type:varchar(100);not null;default:''"`
}
