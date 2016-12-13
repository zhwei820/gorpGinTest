package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/buger/jsonparser"
	"github.com/fatih/structs"
)

// We'll use these two structs to demonstrate encoding and
// decoding of custom types below.
type Response1 struct {
	Page   int
	Fruits []string
}
type Response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

type User struct {
	Id        int64
	Username  string
	Password  string
	Logintime time.Time
}

func jsonTest() {
	user := User{5, "zhangsan", "pwd", time.Now()}
	data1 := structs.Map(user)
	data1["remark"] = "map can be Marshal to json string "
	datas, _ := json.Marshal(data1)
	fmt.Println(string(datas))

	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]interface{}{"apple": 5, "lettuce": " fdssssssssssfsdfsdf7"}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	data := []byte(`{
			"person": {
				"name": {
				"first": "Leonid",
				"last": "Bugaev",
				"fullName": "Leonid Bugaev"
				},
				"github": {
				"handle": "buger",
				"followers": 109
				},
				"avatars": [
				{ "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
				]
			},
			"company": {
				"name": "Acme"
			}
			}`)

	v, _, _, _ := jsonparser.Get(data, "person")

	jsonparser.ArrayEach(v, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ss, _, _, _ := jsonparser.Get(value, "url")
		fmt.Println(string(ss))
	}, "avatars")

	jsonparser.GetInt(data, "person", "github", "followers")
	jsonparser.Get(data, "company")

	if value, _, _, err := jsonparser.Get(data, "company"); err == nil {
		s, _, _, _ := jsonparser.Get(value, "name")
		fmt.Println("tttttttttttttt")
		fmt.Println(string(s))
	}

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		s, _, _, _ := jsonparser.Get(value, "url")
		fmt.Println(string(s))
	}, "person", "avatars")

	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
		return nil
	}, "person", "name")

	// The JSON package can automatically encode your
	// custom data types. It will only include exported
	// fields in the encoded output and will by default
	// use those names as the JSON keys.
	res1D := &Response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	// You can use tags on struct field declarations
	// to customize the encoded JSON key names. Check the
	// definition of `Response2` above to see an example
	// of such tags.
	res2D := &Response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))

	// Now let's look at decoding JSON data into Go
	// values. Here's an example for a generic data
	// structure.
	byt := []byte(`{"num":6.13,"strs":["a","b"], "test":{"apple":5,"lettuce":" fdssssssssssfsdfsdf7"}}`)
	var dat map[string]interface{}

	// Here's the actual decoding, and a check for
	// associated errors.
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	fmt.Println(dat["test"])

	// We can also decode JSON into custom data types.
	// This has the advantages of adding additional
	// type-safety to our programs and eliminating the
	// need for type assertions when accessing the decoded
	// data.
	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := Response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])

}

func main() {
	jsonTest()
}
