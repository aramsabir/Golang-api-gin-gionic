package authentication

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/guizot/go-gin-mongodb/config"
	auth "github.com/guizot/go-gin-mongodb/src/models/auth"
	model "github.com/guizot/go-gin-mongodb/src/models/user"
	"github.com/guizot/go-gin-mongodb/src/security"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	userkey        = "user"
	UserCollection = "user"
	AuthCollection = "auth"
)

type LoginStr struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

func MongoConfig() *mgo.Database {
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
	}
	return db
}

func PostLogin(c *gin.Context) {

	authUser := LoginStr{}
	loginForm := model.User{}
	error := c.Bind(&authUser)

	if error != nil {
		c.JSON(200, gin.H{
			"message": "Error Get Body",
		})
		return
	}
	// loginForm.Prepare()
	email := authUser.Email
	password := authUser.Password

	if email == "" {
		c.JSON(200, gin.H{
			"message": "Error Email or Password should not be empty",
		})
		return
	}
	if password == "" {
		c.JSON(200, gin.H{
			"message": "Error Email or Password should not be empty",
		})
		return
	}

	db := *MongoConfig()
	fmt.Println("MONGO RUNNING: ", db)
	err := db.C(UserCollection).Find(bson.M{}).One(&loginForm)

	if err != nil {
		c.JSON(200, gin.H{
			"message": "User or password is wrong",
		})
		return
	}

	validatePassowrd := security.VerifyPassword(loginForm.Password, password)
	if validatePassowrd != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "User or password is wrong",
			"status":  false,
		})
		return
	}
	tokenString, err := GenerateJWT(loginForm.Email)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Error generating token",
			"status":  false,
		})
		return
	}
	// loginForm.ID

	auth := auth.Auth{}
	auth.Token = tokenString
	auth.UserId = loginForm.ID
	db.C(AuthCollection).Insert(auth)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user", "loginForm": &auth, "token": tokenString, "status": true})
}

var mySigningKey = []byte("HGSh256AramSS")

func GenerateJWT(UserName string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = UserName
	claims["exp"] = time.Now().Add(time.Minute * 300000).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	fmt.Println(err)
	if err != nil {
		fmt.Println("Something when wrong", err.Error())
		return "", err
	}

	return tokenString, nil

}
func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// func AuthRequired(c *gin.Context) {
// 	fmt.Println(c.Request.Header["Authorization"])

// 	session := sessions.Default(c)
// 	user := session.Get(userkey)
// 	if user == nil {
// 		// Abort the request with the appropriate error code
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "status": false})
// 		return
// 	}
// 	// Continue down the chain to handler etc
// 	c.Next()
// }

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

		if token.Valid {
			c.Next()
		}

	} else {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "status": false})
		return
	}

}
