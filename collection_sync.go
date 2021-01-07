package coco

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type collectionIdx struct {
	LastPath string `json:"last_path"`
	Idx      uint64 `json:"idx"`
}

func (c *Collection) collectionSynchronization() {
	last := c.lastPath()
	c.count = c.docCount(last.LastPath)
	open, err := os.OpenFile(last.LastPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 00666)
	if err != nil {
		log.Fatalln(err)
	}
loop:
	for {
		select {
		case data, ex := <-c.writeFile:
			if !ex {
				break loop
			}
			if c.count >= c.conf.MaxSize {
				open.Close()

			}

		}
	}
}

func (c *Collection) lastPath() *collectionIdx {
	idx := path.Join(c.pathRoot, "coco.idx")
	file, err := ioutil.ReadFile(idx)
	if err != nil {
		marshal, err := json.Marshal(&collectionIdx{
			LastPath: path.Join(c.pathRoot, "0.data"),
			Idx:      0,
		})
		if err != nil {
			log.Fatalln(err)
		}
		if err := ioutil.WriteFile(idx, marshal, 00666); err != nil {
			log.Fatalln(err)
		}
		return &collectionIdx{
			LastPath: path.Join(c.pathRoot, "0.data"),
			Idx:      0,
		}
	}
	r := collectionIdx{}
	if err := json.Unmarshal(file, &r); err != nil {
		log.Println(err)
		marshal, err := json.Marshal(&collectionIdx{
			LastPath: path.Join(c.pathRoot, "0.data"),
			Idx:      0,
		})
		if err != nil {
			log.Fatalln(err)
		}
		if err := ioutil.WriteFile(idx, marshal, 00666); err != nil {
			log.Fatalln(err)
		}
		return &collectionIdx{
			LastPath: path.Join(c.pathRoot, "0.data"),
			Idx:      0,
		}
	}

	return &r
}

func (c *Collection) docCount(filepath string) uint64 {
	open, err := os.Open(filepath)
	if err != nil {
		create, err := os.Create(filepath)
		if err != nil {
			log.Fatalln(err)
			return 0
		}
		create.Close()
		return 0
	}
	defer open.Close()
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

func (c *Collection) nextNode(idx *collectionIdx) {

}
