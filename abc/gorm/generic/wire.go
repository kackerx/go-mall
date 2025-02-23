//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

func New() (*App, func(), error) {
	wire.Build(
		NewDB,
		NewUserDao,
		NewPostDao,
		NewUserRepo,
		NewPostRepo,
		NewApp,
	)

	return nil, nil, nil
}
