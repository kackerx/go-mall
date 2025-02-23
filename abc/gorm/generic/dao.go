package main

import (
	"context"

	"github.com/k0kubun/pp/v3"
	"gorm.io/gorm"
)

func (d *Dao[T]) List(ctx context.Context) {
	var (
		a    T
		objs []*T
	)
	err := d.db.Table(a.TableName()).Find(&objs).Error
	if err != nil {
		panic(err)
	}

	pp.Println(objs)
}

func (d *Dao[T]) Get(ctx context.Context, id uint) {
	var (
		a T
	)
	err := d.db.Table(a.TableName()).First(&a, id).Error
	if err != nil {
		panic(err)
	}

	pp.Println(a)
}

type Dao[T Model] struct {
	db *gorm.DB
}

// --- Dao ---

type UserDao = *Dao[User]

func NewUserDao(db *gorm.DB) UserDao {
	return &Dao[User]{db: db}
}

// --- PostDao ---

type PostDao = *Dao[Post]

func NewPostDao(db *gorm.DB) PostDao {
	return &Dao[Post]{db: db}
}
