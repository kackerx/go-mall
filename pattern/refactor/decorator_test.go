package refactor

import (
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestDecorator(t *testing.T) {

}

type MyCache struct {
	Name  string
	cache redis.Cmdable
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func NewUser(name string, age int) *User {
	return &User{Name: name, Age: age}
}

func (u *User) GetName() string {
	return u.Name
}

type Getter struct {
	Name string
}
