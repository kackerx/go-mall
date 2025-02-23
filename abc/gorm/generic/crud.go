package main

import "context"

type CRUD interface {
	List(ctx context.Context)
	Get(ctx context.Context, id uint)
}
