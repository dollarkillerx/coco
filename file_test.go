package coco

import (
	"bufio"
	"fmt"
	json "github.com/json-iterator/go"
	"reflect"

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

func TestFf1(t *testing.T) {
	s := ta{
		Name:        "dollarkiller",
		Age:         22,
		Description: "hp",
	}

	marshal, _ := json.Marshal(s)
	rc := make(map[string]interface{})
	err := json.Unmarshal(marshal, &rc)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rc)

	filter := M{
		"age": M{
			"$<": 30,
		},
		"name": "dollarkiller",
	}

	for k, v := range filter {
		switch r := v.(type) {
		case int, int64, int32:

		case uint64, uint32, uint:

		case float32, float64:

		case M:
			for k2, v2 := range r {
				switch k2 {
				case "$>":
					switch r2 := v2.(type) {
					//case int,int32,int64:
					case int:
						i, ex := rc[k]
						if ex {
							c, ex := i.(int)
							if ex {
								if c > r2 {
									fmt.Println(rc[k])
								}
							}
						}
					case uint, uint32, uint64:

					case float32, float64:

					}
				case "$<":
					fmt.Println(v2, "ini")
					switch r2 := v2.(type) {
					//case int,int32,int64:
					case int:
						i, ex := rc[k]
						if ex {
							c, ex := i.(float64)
							if ex {
								if int(c) < r2 {
									fmt.Println(rc[k])
								}
							}
						}
					case uint, uint32, uint64:

					case float32, float64:

					}
				}
			}
		case string:
			u, ex := rc[k]
			if ex {
				s2, e := u.(string)
				if e {
					if s2 == r {
						fmt.Println(s2)
					}
				}
			}
		}
	}
}

func TestFs(t *testing.T) {
	type ft struct {
		Name string  `json:"name"`
		Age  int     `json:"age"`
		Age2 int64   `json:"age_2"`
		Age3 int32   `json:"age_3"`
		Age4 uint32  `json:"age_4"`
		Age5 uint64  `json:"age_5"`
		Age6 float32 `json:"age_6"`
		Age7 float64 `json:"age_7"`
	}

	p := &ft{
		Name: "sp1",
		Age:  15,
		Age2: -15,
		Age3: 213,
		Age4: 2343,
		Age5: 60,
		Age6: 324,
		Age7: 949.00,
	}

	marshal, _ := json.Marshal(p)
	fmt.Println(string(marshal))

	r := make(map[string]interface{})

	if err := json.Unmarshal(marshal, &r); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("name: ", reflect.TypeOf(r["name"]).String())
	fmt.Println("age: ", reflect.TypeOf(r["age"]).String())
	fmt.Println("age_2: ", reflect.TypeOf(r["age_2"]).String())
	fmt.Println("age_3: ", reflect.TypeOf(r["age_3"]).String())
	fmt.Println("age_4: ", reflect.TypeOf(r["age_4"]).String())
	fmt.Println("age_5: ", reflect.TypeOf(r["age_5"]).String())
	fmt.Println("age_6: ", reflect.TypeOf(r["age_6"]).String())
	fmt.Println("age_7: ", reflect.TypeOf(r["age_7"]).String())
}

func TestFs2(t *testing.T) {
	type ft struct {
		Name string  `json:"name"`
		Age  int     `json:"age"`
		Age2 int64   `json:"age_2"`
		Age3 int32   `json:"age_3"`
		Age4 uint32  `json:"age_4"`
		Age5 uint64  `json:"age_5"`
		Age6 float32 `json:"age_6"`
		Age7 float64 `json:"age_7"`
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)

	results := make([]ft, 0)
	//val := reflect.ValueOf(&results)
	//fmt.Println(val.Kind())
	//
	//elem := val.Elem()
	//fmt.Println(elem.Kind())
	ex1(&results)
}

func ex(s interface{}) {
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Ptr {
		log.Fatalln(typ.Kind())
	}

	//val := reflect.ValueOf(s)

	// 获取数组元素类型
	eleEE := typ.Elem().Elem()
	// 获取数组值
	//valEE := val.Elem()

	// new 对象
	item := reflect.New(eleEE)
	// 对象值
	itemVal := item.Elem()
	name := itemVal.FieldByName("Name")
	age := itemVal.FieldByName("Age")
	name.SetString("dollarkiller")
	age.SetInt(60)

	//valEE.Set(item)
}

func ex1(s interface{}) {
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Ptr {
		log.Fatalln(typ.Kind())
	}

	//val := reflect.ValueOf(s)

	// 获取数组元素类型
	eleEE := typ.Elem().Elem()
	// 获取数组值
	//valEE := val.Elem()
	//fmt.Println(typ.Elem().Elem().Kind())
	//of := reflect.ValueOf(eleEE)

	// new 对象
	//item := reflect.New(eleEE)
	// 对象值
	//itemVal := item.Elem()
	//name := itemVal.FieldByName("Name")
	//age := itemVal.FieldByName("Age")
	//name.SetString("dollarkiller")
	//age.SetInt(60)
	for i := 0; i < eleEE.NumField(); i++ {
		//if of.Field(i).CanInterface() {
		fmt.Println(eleEE.Field(i).Name)
		fmt.Println(eleEE.Field(i).Type.Kind())

		//}
	}

	//valEE.Set(item)
}

type User struct {
	Id   int
	Name string
}

func ChangeSlice(s interface{}) {
	sT := reflect.TypeOf(s)
	if sT.Kind() != reflect.Ptr {
		fmt.Println("参数必须是ptr类型")
		os.Exit(-1)
	}
	sV := reflect.ValueOf(s)
	// 取得数组中元素的类型
	sEE := sT.Elem().Elem()
	// 数组的值
	sVE := sV.Elem()

	// new一个数组中的元素对象
	sON := reflect.New(sEE)
	// 对象的值
	sONE := sON.Elem()
	// 给对象复制
	sONEId := sONE.FieldByName("Id")
	sONEName := sONE.FieldByName("Name")
	sONEId.SetInt(10)
	sONEName.SetString("李四")

	// 创建一个新数组并把元素的值追加进去
	newArr := make([]reflect.Value, 0)
	newArr = append(newArr, sON.Elem())

	// 把原数组的值和新的数组合并
	resArr := reflect.Append(sVE, newArr...)

	// 最终结果给原数组
	sVE.Set(resArr)
}

func TestP2(t *testing.T) {
	users := make([]User, 0)
	ChangeSlice(&users)
	// 这里希望让Users指向ChangeSlice函数中的那个新数组
	fmt.Println(users)
}
