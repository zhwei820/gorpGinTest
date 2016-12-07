package main

import (
	"gorpGinTest/models"

	"github.com/gin-gonic/gin"
)

func main() {
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
