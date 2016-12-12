package endpoint

import (
	"fmt"
	"gorpGinTest/models"
	"log"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"
)

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
		// fmt.Println(q["a"][0])
		fmt.Println("query: " + query)
		fmt.Println(q)
	}

	var users []models.User
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

	var user models.User
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

	var user models.User
	c.Bind(&user)

	if verbose == true {
		fmt.Println(user)
	}

	if user.Name != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&user)
		if err == nil {
			c.JSON(201, user)
		} else {
			models.CheckErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Thea\", \"comment\": \"Queen\" }" http://localhost:8084/api/v1/users
}

// UpdateUser update one user by id
func UpdateUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := true
	id := c.Params.ByName("id")

	var user models.User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)
	if err == nil {
		var json models.User
		c.Bind(&json)

		if verbose == true {
			fmt.Println(json)
		}

		userId, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		user := models.User{
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
				models.CheckErr(err, "Updated failed")
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

	var user models.User
	err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&user)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			models.CheckErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i -X DELETE http://localhost:8084/api/v1/users/1
}
