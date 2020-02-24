package routes

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/guizot/go-gin-mongodb/src/handlers/authentication"
	handle_user "github.com/guizot/go-gin-mongodb/src/handlers/user"
	// jwt "github.com/dgrijalva/jwt-go"
)

type Routes struct {
}

const (
	userkey = "user"
)

var mySigningKey = []byte("HGSh256AramSS")

func (c Routes) StartGin() {
	// r := gin.Default()
	r := gin.Default()
	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))
	api := r.Group("/api")
	r.POST("/login", authentication.PostLogin)
	api.Use(isAuthorized)
	{
		api.GET("/info", handle_user.UserInfo)
		api.GET("/users", handle_user.GetAllUser)
		api.POST("/users", handle_user.CreateUser)
		api.GET("/users/:id", handle_user.GetUser)
		api.PUT("/users/:id", handle_user.UpdateUser)
		api.DELETE("/users/:id", handle_user.DeleteUser)
	}

	r.Run(":8000")
}

func isAuthorized(c *gin.Context) {

	authorization := c.Request.Header.Get("Authorization")

	if c.Request.Header["Authorization"] != nil {
		myToken := authorization[7:len(authorization)]
		token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return mySigningKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "status": false})
			return
		}
		fmt.Println(token.Header["map"])
		if token.Valid {
			c.Next()
		}

	} else {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "status": false})
		return
	}

}

// func AuthRequired(c *gin.Context) {
// 	authorization := c.Request.Header.Get("Authorization")

// 	if authorization != nil{
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "status": false})
// 		return
// 	}
// 	myToken := authorization[7:len(authorization)]

// 	verification, err := jwt.Parse(myToken, func (token, *jwt.Token)(intreface{},error)  {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
// 			return nill, fmt.Errorf("Error")
// 		}
// 		return
// 	})
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "status": false})
// 		return
// 	}
// 	c.Next()

// }
