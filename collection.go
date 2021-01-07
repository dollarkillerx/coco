package coco

import (
	"context"
	"encoding/json"
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
