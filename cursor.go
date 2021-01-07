package coco

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/pkg/errors"
)

type Cursor struct {
	tmpFile string
}

func (c *Cursor) Close() error {
	return os.Remove(c.tmpFile)
}

type itemType struct {
	name     string
	typeKind reflect.Kind
}

func (c *Cursor) All(ctx context.Context, results interface{}) (count int, err error) {
	resultsVal := reflect.ValueOf(results)
	if resultsVal.Kind() != reflect.Ptr {
		return count, errors.New(fmt.Sprintf("results argument must be a pointer to a slice, but was a %s", resultsVal.Kind()))
	}

	sliceVal := resultsVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		return count, errors.New(fmt.Sprintf("results argument must be a pointer to a slice, but was a pointer to %s", sliceVal.Kind()))
	}

	resultsType := reflect.TypeOf(results)
	// 取到单个item类型
	itemElem := resultsType.Elem().Elem()

	// 获取到改struct所有item type
	types := make([]*itemType, 0)
	for i := 0; i < itemElem.NumField(); i++ {
		itemType := itemType{
			name:     itemElem.Field(i).Name,
			typeKind: itemElem.Field(i).Type.Kind(),
		}
		types = append(types, &itemType)
	}

	// 创建返回slice
	newArr := make([]reflect.Value, 0)

	// 读 tmp文件 写Slice
	open, err := os.Open(c.tmpFile)
	if err != nil {
		return count, err
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	for {
		byt, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		itemMap := make(map[string]interface{})
		if err := json.Unmarshal(byt[:len(byt)-1], &itemMap); err != nil {
			log.Println(err)
			continue
		}

		// 创建一个新元素
		newItem := reflect.New(itemElem)
		item := newItem.Elem()
		for k, v := range types {
			tag := itemElem.Field(k).Tag.Get("json")
			searchName := v.name
			if tag != "" {
				searchName = tag
			}
			i, ex := itemMap[searchName]
			if !ex {
				continue
			}

			val := item.FieldByName(v.name)
			switch v.typeKind {
			case reflect.Bool:
				bo, ex := i.(bool)
				if ex {
					val.SetBool(bo)
				}
			case reflect.Int:
				bo, ex := i.(float64)
				if ex {
					val.SetInt(int64(bo))
				}
			case reflect.Int32:
				bo, ex := i.(float64)
				if ex {
					val.SetInt(int64(bo))
				}
			case reflect.Int64:
				bo, ex := i.(float64)
				if ex {
					val.SetInt(int64(bo))
				}
			case reflect.Uint:
				bo, ex := i.(float64)
				if ex {
					val.SetUint(uint64(bo))
				}
			case reflect.Uint32:
				bo, ex := i.(float64)
				if ex {
					val.SetUint(uint64(bo))
				}
			case reflect.Uint64:
				bo, ex := i.(float64)
				if ex {
					val.SetUint(uint64(bo))
				}
			case reflect.Float32:
				bo, ex := i.(float64)
				if ex {
					val.SetFloat(bo)
				}
			case reflect.Float64:
				bo, ex := i.(float64)
				if ex {
					val.SetFloat(bo)
				}
			case reflect.String:
				bo, ex := i.(string)
				if ex {
					val.SetString(bo)
				}
			default:
				continue
			}
		}
		newArr = append(newArr, item)
		count++
	}

	// 回写
	resArr := reflect.Append(sliceVal, newArr...)
	sliceVal.Set(resArr)
	return count, nil
}
