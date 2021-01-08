package coco

import (
	"bufio"
	"log"
	"os"

	json "github.com/json-iterator/go"
)

type search struct {
	tmpPath string      // 临时文件路径
	write   chan []byte // 写文件通道
	close   chan bool
	await   chan bool
}

func newSearch(tmpPath string) *search {
	s := &search{
		tmpPath: tmpPath,
		write:   make(chan []byte),
		close:   make(chan bool),
		await:   make(chan bool),
	}

	go s.file()
	return s
}

func (s *search) Close() {
	close(s.close)
}

func (s *search) Await() {
	for {
		select {
		case <-s.await:
			return
		}
	}
}

func (s *search) file() {
	file, err := os.OpenFile(s.tmpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 00666)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		file.Close()
		close(s.await)
	}()
	writer := bufio.NewWriter(file)
loop:
	for {
		select {
		case <-s.close:
			break loop
		case data := <-s.write:
			if _, err := writer.Write(data); err != nil {
				log.Println(err)
			}
			writer.Flush()
		}
	}
}

func (s *search) Search(file string, filter M) {
	open, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}

	reader := bufio.NewReader(open)
	for {
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		item := make(map[string]interface{})
		if err := json.Unmarshal(bytes[:len(bytes)-1], &item); err != nil {
			log.Println(err)
			continue
		}
		// filter
		if s.filter(item, filter) {
			s.write <- bytes
		}
	}
}

func (s *search) filter(item map[string]interface{}, filter M) bool {
	for k, v := range filter {
		k1, ex := item[k]
		if !ex {
			return false
		}

		switch r1 := v.(type) {
		case string:
			ks, ex := k1.(string)
			if !ex {
				return false
			}

			if ks != r1 {
				return false
			}
		case int:
			ks, ex := k1.(int)
			if !ex {
				return false
			}

			if ks != r1 {
				return false
			}
		case int32:
			ks, ex := k1.(int32)
			if !ex {
				return false
			}

			if ks != r1 {
				return false
			}
		case int64:
			ks, ex := k1.(int64)
			if !ex {
				return false
			}

			if ks != r1 {
				return false
			}
		case uint:
			ks, ex := k1.(uint)
			if !ex {
				return false
			}

			if ks != r1 {
				return false
			}
		case uint32:
			ks, ex := k1.(uint32)
			if !ex {
				return false
			}

			if ks != r1 {
				return false
			}
		case uint64:
			ks, ex := k1.(uint64)
			if !ex {
				return false
			}

			if ks != r1 {
				return false
			}
		case float32:
			ks, ex := k1.(float32)
			if !ex {
				return false
			}

			if ks != r1 {
				return false
			}
		case float64:
			ks, ex := k1.(float64)
			if !ex {
				return false
			}

			if ks != r1 {
				return false
			}
		case M:
			if !s.mFilter(k1, r1) {
				return false
			}
		}
	}

	return true
}

func (s *search) mFilter(item interface{}, filter M) bool {
	it, ex := item.(float64)
	if !ex {
		return false
	}
	for k, v := range filter {
		switch k {
		case "$>":
			if !s.mCompare(">", v, it) {
				return false
			}
		case "$<":
			if !s.mCompare("<", v, it) {
				return false
			}
		default:
			return false
		}
	}

	return true
}

func (s *search) mCompare(symbol string, key interface{}, val float64) bool {
	//find, err := collection.Find(context., M{
	//	"age": M{
	//		"$>": 300,
	//	},
	//})

	// key 查询是输入大值  val: 300
	// val 数据库中大值
	switch k := key.(type) {
	case int:
		if symbol == ">" {
			if float64(k) < val {
				return true
			}
		} else if symbol == "<" {
			if float64(k) > val {
				return true
			}
		}
	case int32:
		if symbol == ">" {
			if float64(k) < val {
				return true
			}
		} else if symbol == "<" {
			if float64(k) > val {
				return true
			}
		}
	case int64:
		if symbol == ">" {
			if float64(k) < val {
				return true
			}
		} else if symbol == "<" {
			if float64(k) > val {
				return true
			}
		}
	case uint:
		if symbol == ">" {
			if float64(k) < val {
				return true
			}
		} else if symbol == "<" {
			if float64(k) > val {
				return true
			}
		}
	case uint32:
		if symbol == ">" {
			if float64(k) < val {
				return true
			}
		} else if symbol == "<" {
			if float64(k) > val {
				return true
			}
		}
	case uint64:
		if symbol == ">" {
			if float64(k) < val {
				return true
			}
		} else if symbol == "<" {
			if float64(k) > val {
				return true
			}
		}
	case float32:
		if symbol == ">" {
			if float64(k) < val {
				return true
			}
		} else if symbol == "<" {
			if float64(k) > val {
				return true
			}
		}
	case float64:
		if symbol == ">" {
			if k < val {
				return true
			}
		} else if symbol == "<" {
			if k > val {
				return true
			}
		}
	default:
		return false
	}

	return false
}
