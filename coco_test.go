package coco

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

type ta2 struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Description string `json:"description"`
}

func TestWrite(t *testing.T) {
	client := NewClient(NewDefaultConfig("../out"))
	collection, err := client.Database("test1").Collection("test1")
	if err != nil {
		log.Fatalln(err)
	}

	rs := make([]interface{}, 0)
	for i := 0; i < 5000000; i++ {
		r := ta2{
			Name:        fmt.Sprintf("scp-%d", i),
			Age:         rand.Intn(600),
			Description: fmt.Sprintf("scp-NB-%d", i),
		}
		rs = append(rs, r)
	}

	now := time.Now()
	err = collection.InsertMany(context.TODO(), rs)
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Millisecond * 1)
	since := time.Since(now)
	fmt.Println(since.Milliseconds()) // 7325ms  =>  500W
}

func TestRead(t *testing.T) {
	now := time.Now()
	client := NewClient(NewDefaultConfig("../out"))
	collection, err := client.Database("test1").Collection("test1")
	if err != nil {
		log.Fatalln(err)
	}

	find, err := collection.Find(context.TODO(), M{
		//"name": "dollarkiller",
		"age": M{
			"$<": 30,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(find)
	//defer find.Close()
	fmt.Println("s1: ", time.Since(now).Milliseconds())

	rs := make([]ta2, 0)
	all, err := find.All(context.TODO(), &rs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(all)
	fmt.Println("s2: ", time.Since(now).Milliseconds())

	time.Sleep(time.Second)
}
