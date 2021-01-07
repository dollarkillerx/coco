package coco

import (
	"bufio"
	"encoding/json"
	"fmt"
	//"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

type ta struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Description string `json:"description"`
}

func TestFIleOut(t *testing.T) {
	open, err := os.Create("./out/f.doc")
	if err != nil {
		log.Fatalln(err)
	}
	defer open.Close()

	for i := 0; i < 100; i++ {
		rand.Seed(time.Now().UnixNano())
		r := ta{
			Name:        fmt.Sprintf("scp-%d", i),
			Age:         rand.Intn(600),
			Description: fmt.Sprintf("scp-NB-%d", i),
		}
		marshal, err := json.Marshal(&r)
		if err != nil {
			log.Println(err)
			continue
		}

		open.Write(append(marshal, +'\n'))
	}
}

//func TestFIleBsonOut(t *testing.T) {
//	open, err := os.Create("./out/fb.doc")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer open.Close()
//
//	for i := 0; i < 100; i++ {
//		rand.Seed(time.Now().UnixNano())
//		r := ta{
//			Name:        fmt.Sprintf("scp-%d", i),
//			Age:         rand.Intn(600),
//			Description: fmt.Sprintf("scp-NB-%d", i),
//		}
//		marshal, err := bson.Marshal(&r)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//
//		open.Write(append(marshal, +'\n'))
//	}
//}
//
//func TestRBsoni(t *testing.T) {
//	open, err := os.Open("./out/fb.doc")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer open.Close()
//
//	reader := bufio.NewReader(open)
//
//	for {
//		bytes, err := reader.ReadBytes('\n')
//		if err != nil {
//			break
//		}
//		r := ta{}
//		if err := bson.Unmarshal(bytes[:len(bytes)-1], &r); err != nil {
//			log.Fatalln(err)
//		}
//
//		fmt.Println(r)
//	}
//}


func TestRi(t *testing.T) {
	open, err := os.Open("./out/f.doc")
	if err != nil {
		log.Fatalln(err)
	}
	defer open.Close()

	reader := bufio.NewReader(open)

	for {
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		r := ta{}
		if err := json.Unmarshal(bytes[:len(bytes)-1], &r); err != nil {
			log.Fatalln(err)
		}

		fmt.Println(r)
	}
}

func TestFileCount(t *testing.T) {
	count := fileCount("./out/f.doc")
	fmt.Println(count)

	count = fileCount("./out/f2.doc")
	fmt.Println(count)
}

func fileCount(filepath string) uint64 {
	open, err := os.Open(filepath)
	if err != nil {
		create, err := os.Create(filepath)
		if err != nil {
			log.Fatalln(err)
		}
		create.Close()
		return 0
	}
	reader := bufio.NewReader(open)
	var count uint64
	for {
		_, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		count++
	}
	return count
}
