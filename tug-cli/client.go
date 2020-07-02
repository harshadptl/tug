package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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


		print("Enter c to continue or f/F to flush logs or q/Q to quit:")
		inp, _ := reader.ReadString('\n')
		inp = strings.ToLower(inp)


		if inp == "c\n" {

			subscriptions, err := cl.PubSubNumSub("tug").Result()
			if s, ok := subscriptions["tug"]; !ok || s < 1 {
				print("no checkpoint...\n")
				continue
			}

			err = cl.Publish("tug", "go").Err()
			if err != nil {
				print("redis error: ", err.Error(), "\n")
			}
		} else if inp == "f\n" {

			err := cl.Del("tug").Err()
			if err != nil {
				print("error flushing logs: ", err.Error(), "\n")
			}
		} else if inp == "q\n" {
			goto done
		} else {

			val, err := cl.XRange("tug", "-", "+").Result()
			if err != nil {
				print("redis error: ", err.Error())
				print("\n")
			}

			ids := []string{}

			print()
			for i := range val {
				msg := val[i]
				print("\n\n\n")
				print("Printing logs-----------------------------------------------------------------------------\n")
				{
					ids := strings.Split(msg.ID, "-")
					timestamp := ids[0]
					i, err := strconv.ParseInt(timestamp, 10, 64)
					if err != nil {

					}
					print(time.Unix(0, i*1000000).String())
				}
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
		}

	}
	done:
}
