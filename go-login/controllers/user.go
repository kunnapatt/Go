package controllers

import (
	"fmt"
	"net/http"

	"test/go-login/orm"

	"github.com/gin-gonic/gin"
)

func ReadAll(c *gin.Context) {

	var users []orm.User

	orm.Db.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "User Read Success",
		"users": users,
	})
}

func Profile(c *gin.Context) {

	userId := c.MustGet("userId").(float64)
	fmt.Println(userId)
	var user orm.User


	orm.Db.First(&user, userId)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "User Read Success",
		"users": user,
	})
}