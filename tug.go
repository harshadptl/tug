package tug

import (
	"reflect"
	"fmt"
	"time"
	"sync"
	"github.com/go-redis/redis/v7"
)

const (
	DbNumber = 15
)
var once sync.Once

type Tug struct {
	name string
	client *redis.Client
}

func NewTug(name, host, password string) *Tug {

	cl := redis.NewClient(&redis.Options{
		Addr: host,
		Password: password,
		DB: DbNumber,
	})

	return &Tug{
		name: name,
		client: cl,
	}
}

func (t *Tug) Pause(vars ...interface{}) {

	id := fmt.Sprintf("%s: %s", t.name, time.Now().String())

	xa := &redis.XAddArgs{
		Stream: "tug",
		ID:	id,
		Values: map[string]interface{}{},
	}

	for i := range vars {
		typ := reflect.ValueOf(vars[i]).Type()
		key := fmt.Sprintf("%d:%s", i, reflect.ValueOf(typ))

		xa.Values[key] = vars[i]
	}

	_, err := t.client.XAdd(xa).Result()
	if err != nil {
		panic("tug log failed")
	}

	block := t.client.Subscribe("tug").Channel()
	<-block

}



