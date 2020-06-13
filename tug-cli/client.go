package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v7"
)

const (
	DbNumber = 15
)

func main() {
	args := os.Args

	if len(args) != 3 {
		print("Usage: tug-cli [host] [password]\n")
		return
	}

	host := os.Args[1]
	password := os.Args[2]

	cl := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       DbNumber,
	})

	reader := bufio.NewReader(os.Stdin)

	for true {
		val, err := cl.XRange("tug", "-", "+").Result()
		if err != nil {
			print("redis error: ", err.Error())
		}

		ids := []string{}

		print()
		for i := range val {
			msg := val[i]
			print("-------------------------------------------------")
			print(msg.ID)
			ids = append(ids, msg.ID)

			for k, v := range msg.Values {
				fmt.Printf("|%s\t\t|%v", k, v)
			}

		}
		if len(ids) == 0 {
			print("no logs...")
		}

		print("Press c/C to continue or f/F to flush logs or q/Q to quit:")
		inp, _ := reader.ReadString('\n')
		inp = strings.ToLower(inp)

		switch inp {
		case "c":
			err := cl.Publish("tug", "go").Err()
			if err != nil {
				print("redis error: ", err.Error())
			}
		case "f":
			err := cl.XDel("tug", ids...).Err()
			if err != nil {
				print("error flushing logs: ", err.Error())
			}
		case "q":
			return
		}


	}
}
