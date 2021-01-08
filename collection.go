package coco

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dollarkillerx/async_utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

func (c *Collection) InsertMany(ctx context.Context, documents []interface{}) error {
	over := make(chan bool)
	poolFunc := async_utils.NewPoolFunc(int(c.conf.MaxWritePool), func() {
		close(over)
	})

	var err error
	for k := range documents {
		idx := k
		poolFunc.Send(func() {
			marshal, err := json.Marshal(documents[idx])
			if err != nil {
				err = err
				log.Println(err)
				return
			}

			c.writeFile <- marshal
		})
	}
	poolFunc.Over()

	<-over
	return err
}

func (c *Collection) Find(ctx context.Context, filter M) (*Cursor, error) {
	// 先查询写临时表
	paths, err := ioutil.ReadDir(c.pathRoot)
	if err != nil {
		return nil, errors.New("Collection 404")
	}

	over := make(chan bool)
	poolFunc := async_utils.NewPoolFunc(int(c.conf.MaxReadPool), func() {
		over <- true
	})

	tmpFile := path.Join(c.pathRoot, fmt.Sprintf("%s.tmp", uuid.New().String()))
	s := newSearch(tmpFile)
	for k := range paths {
		idx := k
		if !paths[idx].IsDir() && strings.Index(paths[idx].Name(), ".data") != -1 {
			poolFunc.Send(func() {
				s.Search(path.Join(c.pathRoot, paths[idx].Name()), filter)
			})
		}
	}
	poolFunc.Over()

	<-over
	s.Close()
	s.Await()
	dir, err := ioutil.ReadFile(tmpFile)
	if err == nil {
		log.Println("Dir: ",string(dir))
	}
	return &Cursor{tmpFile: tmpFile}, nil
}
