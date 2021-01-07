package coco

import (
	"context"
	"encoding/json"
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
	client := NewClient(NewDefaultConfig("./out"))
	collection, err := client.Database("test1").Collection("test1")
	if err != nil {
		log.Fatalln(err)
	}

	rs := make([]interface{}, 0)
	for i := 0; i < 20300; i++ {
		r := ta2{
			Name:        fmt.Sprintf("scp-%d", i),
			Age:         rand.Intn(600),
			Description: fmt.Sprintf("scp-NB-%d", i),
		}
		rs = append(rs, r)
	}
	err = collection.InsertMany(context.TODO(), rs)
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Second * 10)
}

func TestRead(t *testing.T) {
	client := NewClient(NewDefaultConfig("./out"))
	collection, err := client.Database("test1").Collection("test1")
	if err != nil {
		log.Fatalln(err)
	}

	find, err := collection.Find(context.TODO(), M{
		"age": M{
			"$<": 300,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(find)
	defer find.Close()

	rs := make([]ta2, 0)
	all, err := find.All(context.TODO(), &rs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(all)

	marshal, err := json.Marshal(rs)
	if err == nil {
		log.Println(string(marshal))
	}

	time.Sleep(time.Second)
}
