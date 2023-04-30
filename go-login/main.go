package main

import (
	"fmt"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	Controller "test/go-login/controllers"
	"test/go-login/orm"
	"test/go-login/middlewares"
)


var (
	USER="root"
	PASS=""
	localhost="127.0.0.1"

)

func main() {

	initENV()
	orm.InitDB()

	fmt.Println("starting server")

	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/register", Controller.Register)
	r.POST("/login", Controller.Login)

	authorized := r.Group("/users", middlewares.JWTAuthen())
	authorized.GET("/readall", Controller.ReadAll)
	authorized.GET("/profile", Controller.Profile)

	// r.GET("/users/readall", Controller.ReadAll)
	r.Run(":8080")
}

func initENV() {
	err := godotenv.Load(".env")

	if (err != nil) {
		fmt.Println("Error loading .env file")
	}
}