package middlewares

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"

	"test/go-login/orm"

)

func JWTAuthen() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set example variable
		claims, check := checkJWT(c)

		if (check) {
			c.Set("userId", claims["userId"])
		} else {

		}

		// before request

		c.Next()

	}
}

func checkJWT(c *gin.Context) (jwt.MapClaims, bool){
	hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	header := c.Request.Header.Get("Authorization")
	tokenString := strings.Replace(header, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		fmt.Println(claims["userId"], claims["nbf"])
		var users []orm.User

		orm.Db.Find(&users)

		return claims, true
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": err.Error(),
		})

		return jwt.MapClaims{}, false
	}
}
