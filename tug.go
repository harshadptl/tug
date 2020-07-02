package tug

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"reflect"
)

const (
	DbNumber = 15
)

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


	xa := &redis.XAddArgs{
		Stream: "tug",
		ID:	"*",
		Values: map[string]interface{}{},
	}

	xa.Values["appName"] = t.name

	for i := range vars {
		typ := reflect.ValueOf(vars[i]).Type()
		key := fmt.Sprintf("%d:%s", i, reflect.ValueOf(typ))
		val := fmt.Sprintf("%#v", vars[i])

		xa.Values[key] = val
	}
	_, err := t.client.XAdd(xa).Result()
	if err != nil {
		panic("tug log failed: " + err.Error())

	}

	sub := t.client.Subscribe("tug")
	block := sub.Channel()
	<- block
	sub.Unsubscribe("tug")

}

func (t *Tug) Print(vars ...interface{}) {

	xa := &redis.XAddArgs{
		Stream: "tug",
		ID:	"*",
		Values: map[string]interface{}{},
	}

	xa.Values["appName"] = t.name

	for i := range vars {
		typ := reflect.ValueOf(vars[i]).Type()
		key := fmt.Sprintf("%d:%s", i, reflect.ValueOf(typ))
		val := fmt.Sprintf("%#v", vars[i])

		xa.Values[key] = val
	}
	_, err := t.client.XAdd(xa).Result()
	if err != nil {
		panic("tug log failed: " + err.Error())

	}
}

func (t *Tug) Flush() {
	t.client.Del("tug")
}


