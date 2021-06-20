package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)


func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	key := "ticket_count"
	client.Set(context.Background(),key, "5", 0).Err()
	val, _ := client.Get(context.Background(), key).Result()
	fmt.Println("current ticket_count key val: ", val)

	getTicket(client, key)
}

func runTx(key string) func(tx *redis.Tx) error {
	txf := func(tx *redis.Tx) error {
		n, err := tx.Get(context.Background(),key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		if n == 0 {
			return errors.New("n ==0 ")
		}

		// actual operation (local in optimistic lock)
		n = n - 1

		// runs only if the watched keys remain unchanged
		_, err = tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
			// pipe handles the error case
			pipe.Set(context.Background(), key, n, 0)
			return nil
		})
		return err
	}
	return txf
}

func getTicket(client *redis.Client, key string) {
	err := client.Watch(context.Background(), runTx(key), key)
	if err != nil{
		log.Fatal(err)

	}
}

