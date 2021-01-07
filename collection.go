package coco

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/dollarkillerx/async_utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *Collection) InsertMany(ctx context.Context, documents []interface{}) error {
	for _, v := range documents {
		marshal, err := json.Marshal(v)
		if err != nil {
			return err
		}

		c.writeFile <- marshal
	}

	return nil
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
	return &Cursor{tmpFile: tmpFile}, nil
}
