// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

// Injectors from wire.go:

func New() (*App, func(), error) {
	gormDB := NewDB()
	dao := NewUserDao(gormDB)
	userRepo := NewUserRepo(dao)
	mainDao := NewPostDao(gormDB)
	postRepo := NewPostRepo(mainDao)
	app, cleanup, err := NewApp(userRepo, postRepo)
	if err != nil {
		return nil, nil, err
	}
	return app, func() {
		cleanup()
	}, nil
}
