package main

import (
	"encoding/json"
	"fmt"
	"gorpGinTest/models"
	"time"

	"github.com/gin-gonic/gin"
)
import "github.com/fatih/structs"

func main() {
	jsonTest()
	return

	r := gin.Default()

	r.Use(models.Database("root:spwx@/todolist"))
	r.Use(models.RedisPool("redis/localhost:6379/1", "", 10))

	v1 := r.Group("api/v1")
	{
		v1.GET("/users", models.GetUsers)
		v1.GET("/users/:id", models.GetUser)
		v1.POST("/users", models.PostUser)
		v1.PUT("/users/:id", models.UpdateUser)
		v1.DELETE("/users/:id", models.DeleteUser)
		// v1.OPTIONS("/users", Options)     // POST
		// v1.OPTIONS("/users/:id", Options) // PUT, DELETE
	}
	r.Run(":8084")
}

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
	data := structs.Map(user)
	data["remark"] = "map can be Marshal to json string "
	datas, _ := json.Marshal(data)
	fmt.Println(string(datas))

	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]interface{}{"apple": 5, "lettuce": " fdssssssssssfsdfsdf7"}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

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
