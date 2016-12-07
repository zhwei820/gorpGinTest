package models

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

/**
Search for XXX to fix fields mapping in Update handler, mandatory fields
or remove sqlite tricks

 vim search and replace cmd to customize struct, handler and instances
  :%s/User/NewStruct/g
  :%s/user/newinst/g

**/

// XXX custom struct name and fields
// User db and json type
type User struct {
	Id      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Email   string `db:"email" json:"mail"`
	Status  string `db:"status" json:"status"`
	Comment string `db:"comment, size:16384" json:"comment"`
	Pass    string `db:"pass" json:"pass"`
	Created string `db:"created" json:"created"` // or int64
	Updated string `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

// PreInsert set created an updated time before insert in db
func (a *User) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now().Format("2006-01-02 15:04:05") // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

// PreUpdate set updated time before insert in db
func (a *User) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

// REST handlers

// GetUsers return all users filtered by URL query
func GetUsers(c *gin.Context) {
	dbmap := c.Keys["DBmap"].(*gorp.DbMap)
	pool := c.Keys["Pool"].(*redis.Pool)
	fmt.Println(pool)

	con := pool.Get()
	defer con.Close()

	key := "user1"
	user, err1 := redis.String(con.Do("GET", key))
	if err1 != nil {
		log.Println("user check err1or", err1)

	}
	fmt.Println(user)

	verbose := true
	query := "SELECT * FROM user"

	q := c.Request.URL.Query()
	// query = query + ParseQuery(q)
	if verbose == true {
		fmt.Println(q["a"][0])
		fmt.Println("query: " + query)
	}

	var users []User
	_, err := dbmap.Select(&users, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(len(users)))
		c.JSON(200, users)
	} else {
		c.JSON(404, gin.H{"error": "no user(s) into the table"})
	}

	// curl -i http://localhost:8084/api/v1/users
}

// GetUser return one user by id
func GetUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, user)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i http://localhost:8084/api/v1/users/1
}

// PostUser create and return one user
func PostUser(c *gin.Context) {
	dbmap := c.Keys["DBmap"].(*gorp.DbMap)
	verbose := true

	var user User
	c.Bind(&user)

	if verbose == true {
		fmt.Println(user)
	}

	if user.Name != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&user)
		if err == nil {
			c.JSON(201, user)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}

	// type User struct {
	// 	Id      int64     `db:"id" json:"id"`
	// 	Name    string    `db:"name" json:"name"`
	// 	Email   string    `db:"email" json:"mail"`
	// 	Status  string    `db:"status" json:"status"`
	// 	Comment string    `db:"comment, size:16384" json:"comment"`
	// 	Pass    string    `db:"pass" json:"pass"`
	// 	Created time.Time `db:"created" json:"created"` // or int64
	// 	Updated time.Time `db:"updated" json:"updated"`
	// }
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Thea\", \"comment\": \"Queen\" }" http://localhost:8084/api/v1/users
}

// UpdateUser update one user by id
func UpdateUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := true
	id := c.Params.ByName("id")

	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)
	if err == nil {
		var json User
		c.Bind(&json)

		if verbose == true {
			fmt.Println(json)
		}

		userId, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		user := User{
			Id:      userId,
			Pass:    json.Pass,
			Name:    json.Name,
			Email:   json.Email,
			Status:  json.Status,
			Comment: json.Comment,
			Created: user.Created, //user read from previous select
		}

		if user.Name != "" { // XXX Check mandatory fields
			_, err = dbmap.Update(&user)
			if err == nil {
				c.JSON(200, user)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"name\": \"Merlyn\", \"comment\": \"Merlyn\" }" http://localhost:8084/api/v1/users/1
}

// DeleteUser delete one user by id
func DeleteUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var user User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&user)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i -X DELETE http://localhost:8084/api/v1/users/1
}
