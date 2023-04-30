package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"test/go-login/orm"
)

var (
	hmacSampleSecret []byte
)

type RegisterBody struct {
	Username string `json: "username" binding: "required"`
	Password string `json: "password" binding: "required"`
	Fullname string `json: fullname binding: "required"`
	Avatar string `json: avatar binding: "required"`
}

type LoginBody struct {
	Username string `json: "username" binding: "required"`
	Password string `json: "password" binding: "required"`
}

func Register(c *gin.Context) {
	var json RegisterBody

	orm.InitDB()

	if err := c.ShouldBindJSON(&json) ; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	///// check user exist
	_, exist := checkUserExist(&json)

	if (exist) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "User Exists",
		})
		return 
	}

	encryptPassword, _ := bcrypt.GenerateFromPassword([] byte(json.Password), 10)

	user := orm.User{
		Username: json.Username,
		Password: string(encryptPassword),
		Fullname: json.Fullname,
		Avatar: json.Avatar,
	}

	orm.Db.Create(&user)

	if (user.ID > 0) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"userId": user.ID,
			"message": "User Created",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"message": "User Failed",
		})
	}


}

func checkUserExist(json *RegisterBody) (orm.User, bool){
	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	fmt.Println(userExist.ID)
	if (userExist.ID > 0) {
		return userExist, true
	}
	return userExist, false
}

func Login(c *gin.Context) {
	var json RegisterBody

	orm.InitDB()

	if err := c.ShouldBindJSON(&json) ; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userexist, exist := checkUserExist(&json)

	if (exist) {
		err := bcrypt.CompareHashAndPassword([]byte(userexist.Password), []byte(json.Password))
		///// success
		if (err == nil) {

			hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userId": userexist.ID,
				"exp": time.Now().Add(time.Minute + 1).Unix(),
				// "nbf": time.Date(2023, 10, time.UTC).Unix(),
			})

			tokenString, err := token.SignedString(hmacSampleSecret)
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"messsage": "Login successed.",
				"token": tokenString,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusBadRequest,
				"messsage": "Login failed",
			})
		}
	}

	return

}
