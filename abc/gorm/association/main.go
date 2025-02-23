package main

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/kackerx/go-mall/dal/dao"
)

var (
	db = dao.DBMaster()
)

func init() {
	if err := db.AutoMigrate(&User{}, &Company{}, &Card{}, &Post{}, &Language{}); err != nil {
		panic(err)
	}
}

func addUser() error {
	user := User{
		Name: "jinzhu",
		Company: Company{
			Name: "jinzhu company",
		},
	}
	return db.Create(&user).Error
}

func main() {
	// if err := addUser(); err != nil {
	// 	panic(err)
	// }

	// belongsTo()
	// hasOne()
	// hasMany()
	// Joins()
	// delAssociations()
	delManyToManyAssociations()
}

func delManyToManyAssociations() {
	var user User = User{
		ID:        3,
		Languages: []Language{},
	}
	err := db.Model(&user).
		// Association("Languages").Delete(&Language{ID: 1}) // 根据id删除关联关系
		Association("Languages").Clear() // 清空关联关系
	if err != nil {
		panic(err)
	}

	err = db.Model(&user).Preload("Languages").First(&user).Error
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", user)
}

func delAssociations() {
	var user User = User{
		ID:    3,
		Posts: []Post{{ID: 1}},
	}
	err := db.
		Select("Posts"). // select删除和主表关联的从表数据
		Delete(&user).Error
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", user)
}

func Joins() {
	type Tmp struct {
		ID        int
		Name      string
		PostTitle string `gorm:"column:Posts__title"`
	}
	var (
		user User
		// posts Post
		tmp Tmp
	)
	err := db.Model(User{}).
		Joins("Posts").
		Where("users.id = ?", 3).
		Select("name").
		Scan(&tmp).Error
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", user)
	fmt.Printf("%+v\n", tmp)
}

func hasMany() {
	var user User
	db.
		// Model(User{}).
		// Preload("Posts").
		// Preload(clause.Associations). // all Association
		// Preload("Posts", "id <> ?", 1).
		Preload("Posts", func(db *gorm.DB) *gorm.DB {
			return db.Where("id <> ?", 1)
		}).
		Where("id = ?", 3).
		First(&user)

	fmt.Printf("%+v\n", user)
}

func hasOne() {
	var user User
	db.Model(User{}).
		Preload("Card").
		Preload("Manager").
		Where("id = ?", 3).
		First(&user)

	fmt.Printf("%+v\n", user)
}

func belongsTo() {
	var (
		user    = User{}
		company = Company{}
	)
	err := db.Model(User{}).
		Preload("Company").
		First(&user).Error
	// Find(&company).
	if err != nil {
		panic(err)
	}

	fmt.Println(user, company)
}
