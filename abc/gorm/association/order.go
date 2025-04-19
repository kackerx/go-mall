package main

type Company struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"type:varchar(100);not null;default:''"`

	CompanyExtra *CompanyExtra `gorm:"foreignkey:CompanyID"`
}

type CompanyExtra struct {
	ID        int    `gorm:"primary_key"`
	Name      string `gorm:"type:varchar(100);not null;default:''"`
	CompanyID int    `gorm:"index";`
}
