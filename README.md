# coco
Built-in document database, Large-scale data warehouse

### Use 
- insert
```go
type demo struct {
    Name        string `json:"name"`
	Age         int    `json:"age"`
	Description string `json:"description"`
}

client := NewClient(NewDefaultConfig("./out"))
collection, err := client.Database("test1").Collection("test1")
if err != nil {
    log.Fatalln(err)
}

rs := make([]interface{}, 0)
for i := 0; i < 60030; i++ {
    r := demo{
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
```
- search
```go
    client := NewClient(NewDefaultConfig("./out"))
	collection, err := client.Database("test1").Collection("test1")
	if err != nil {
		log.Fatalln(err)
	}

	find, err := collection.Find(context.TODO(), M{
        "name": "scp",        
		"age": M{
			"$<": 300,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(find)
	defer find.Close()

	rs := make([]demo, 0)
	all, err := find.All(context.TODO(), &rs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(all)

	marshal, err := json.Marshal(rs)
	if err == nil {
		log.Println(string(marshal))
	}
```

### BateV1 Performance
- Insert (7500 HHD)   30W/s
- Search (7500 HHD No Index )   500W data full table scan time 3s