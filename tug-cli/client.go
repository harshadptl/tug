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

	if len(args) != 2 && len(args) != 3 {
		print("Usage: tug-cli [host] [password]\n")
		return
	}

	host := os.Args[1]
	password := ""
	if len(args) == 3 {
		password = os.Args[2]
	}
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
			print("\n\n\n")
			print("Printing logs-------------------------------------------------")
			print(msg.ID)
			ids = append(ids, msg.ID)

			for k, v := range msg.Values {
				print("\n")
				fmt.Printf("|%s\t\t|%v", k, v)
			}

		}
		print("\n\n\n")
		if len(ids) == 0 {
			print("no logs... \n")
		}

		print("Press c to continue or f/F to flush logs or q/Q to quit:")
		inp, _ := reader.ReadString('\n')
		inp = strings.ToLower(inp)


		if inp == "c\n" {

			cl.PubSubChannels("tug").Result()

			err := cl.Publish("tug", "go").Err()
			if err != nil {
				print("redis error: ", err.Error(), "\n")
			}
		} else if inp == "f\n" {

			if len(ids) == 0 {
				continue
			}

			err := cl.Del("tug").Err()
			if err != nil {
				print("error flushing logs: ", err.Error(), "\n")
			}
		} else if inp == "q\n" {
			goto done
		}

	}
	done:
}
