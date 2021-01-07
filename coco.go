package coco

import (
	"path"

	"github.com/dollarkillerx/coco/utils"
)

type Client struct {
	conf *Conf
}

func NewClient(conf *Conf) *Client {
	return &Client{
		conf: conf,
	}
}

type Database struct {
	client   *Client
	database string
}

func (c *Client) Database(database string) *Database {
	return &Database{
		client:   c,
		database: database,
	}
}

type Collection struct {
	database   *Database
	conf       *Conf
	collection string
	pathRoot   string

	count     uint64 // 当前文件count
	writeFile chan []byte
}

func (d *Database) Collection(collection string) (*Collection, error) {
	path := path.Join(d.client.conf.Path, d.database, collection)
	err := utils.DirPing(path)

	c := &Collection{
		database:   d,
		collection: collection,
		conf:       d.client.conf,
		pathRoot:   path,
		writeFile:  make(chan []byte, 100),
	}
	go c.collectionSynchronization()
	return c, err
}

