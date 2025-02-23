package main

import (
	"context"

	"gorm.io/gorm"

	"github.com/kackerx/go-mall/dal/dao"
)

func NewDB() *gorm.DB {
	return dao.DBMaster()
}

type UserRepo struct {
	dao CRUD
}

func NewUserRepo(dao UserDao) *UserRepo {
	return &UserRepo{dao: dao}
}

var (
	db  = dao.DBMaster()
	ctx = context.Background()
)

func init() {
	if err := db.AutoMigrate(&User{}, &Card{}, &Post{}, &Language{}, &Company{}); err != nil {
		panic(err)
	}
}

type PostRepo struct {
	dao CRUD
}

func (p *PostRepo) Get(ctx context.Context, id uint) {
	p.dao.Get(ctx, id)
}

func NewPostRepo(dao PostDao) *PostRepo {
	return &PostRepo{dao: dao}
}

func main() {
	app, _, err := New()
	if err != nil {
		panic(err)
	}

	app.userRepo.ListUser(ctx)
	app.userRepo.Get(ctx, 3)

	app.postRepo.Get(ctx, 3)
}

func (u *UserRepo) ListUser(ctx context.Context) {
	u.dao.List(ctx)
}

func (u *UserRepo) Get(ctx context.Context, id uint) {
	u.dao.Get(ctx, id)
}

type App struct {
	userRepo *UserRepo
	postRepo *PostRepo
}

func NewApp(userRepo *UserRepo, postRepo *PostRepo) (*App, func(), error) {
	return &App{
		userRepo: userRepo,
		postRepo: postRepo,
	}, nil, nil
}
