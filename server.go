package main

import (
	"gorpGinTest/endpoint"
	"gorpGinTest/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()

	r.Use(models.Database("root:spwx@/todolist"))
	r.Use(models.RedisPool("redis://:@localhost:6379/1", "", 10))

	v1 := r.Group("api/v1")
	{
		v1.GET("/users", endpoint.GetUsers)
		v1.GET("/users/:id", endpoint.GetUser)
		v1.POST("/users", endpoint.PostUser)
		v1.PUT("/users/:id", endpoint.UpdateUser)
		v1.DELETE("/users/:id", endpoint.DeleteUser)
		// v1.OPTIONS("/users", Options)     // POST
		// v1.OPTIONS("/users/:id", Options) // PUT, DELETE
	}
	r.Run(":8084")
}
