package user

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guizot/go-gin-mongodb/config"
	auth "github.com/guizot/go-gin-mongodb/src/models/auth"
	model "github.com/guizot/go-gin-mongodb/src/models/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Static Collection

const (
	UserCollection = "user"
	AuthCollection = "auth"
)

// Get DB from Mongo Config
func MongoConfig() *mgo.Database {
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
	}
	return db
}

// Get All User Endpoint
func GetAllUser(c *gin.Context) {
	UserInfo := FindUserByToken(c)

	if UserInfo.Role != "admin" {
		if UserInfo.Role != "developer" {
			c.JSON(200, gin.H{
				"message": "You are not authorize",
				"status":  false,
			})
			return
		}
	}

	db := *MongoConfig()
	fmt.Println("MONGO RUNNING: ", db)

	// users := model.Users{}
	pipeline := []bson.M{
		bson.M{"$lookup": bson.M{"from": "user", "localField": "creator", "foreignField": "_id", "as": "owner"}},
		bson.M{"$unwind": bson.M{"path": "$owner", "preserveNullAndEmptyArrays": true}},
	}
	resp := []bson.M{}
	pipe := db.C(UserCollection).Pipe(pipeline)
	err := pipe.All(&resp)

	if err != nil {
		c.JSON(200, gin.H{
			"message": "Error Get All User",
			"status":  false,
		})
		return
	}

	c.JSON(200, gin.H{
		"user": &resp, "status": true,
	})
	return
}

// Get User Endpoint
func GetUser(c *gin.Context) {
	db := *MongoConfig()
	fmt.Println("MONGO RUNNING: ", db)

	id := c.Param("id") // Get Param
	// idParse, errParse := strconv.Atoi(id) // Convert String to Int
	// if errParse != nil {
	// 	c.JSON(200, gin.H{
	// 		"message": "Error Parse Param",
	// 	})
	// 	return
	// }

	if id == "" {
		c.JSON(200, gin.H{
			"message": "User not valid",
			"status":  false,
		})
		return
	}

	pipeline := []bson.M{
		bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "creator", "foreignField": "_id", "as": "owner"}},
		bson.M{"$unwind": bson.M{"path": "$owner", "preserveNullAndEmptyArrays": true}},
	}

	resp := bson.M{}
	pipe := db.C(UserCollection).Pipe(pipeline)
	err := pipe.One(&resp)

	if err != nil {
		c.JSON(200, gin.H{
			"message": "Error Get User",
		})
		return
	}

	c.JSON(200, gin.H{
		"user":   &resp,
		"status": true,
	})
}

// Create User Endpoint
func CreateUser(c *gin.Context) {
	UserInfo := FindUserByToken(c)

	if UserInfo.Role != "admin" {
		if UserInfo.Role != "developer" {
			c.JSON(200, gin.H{
				"message": "You are not authorize",
				"status":  false,
			})
			return
		}
	}
	db := *MongoConfig()
	fmt.Println("MONGO RUNNING: ", db)

	user := model.User{}
	err := c.Bind(&user)

	if err != nil {

		c.JSON(200, gin.H{
			"message": "Error Get Body",
			"status":  false,
		})
		return
	}

	user.Prepare()

	if user.Name == "" {
		c.JSON(200, gin.H{
			"message": "Error Name required", "status": false,
		})
		return
	}
	if user.Email == "" {
		c.JSON(200, gin.H{
			"message": "Error Email required", "status": false,
		})
		return
	}
	if user.Password == "" {
		c.JSON(200, gin.H{
			"message": "Error Password required", "status": false,
		})
		return
	}
	if user.Role == "" {
		c.JSON(200, gin.H{
			"message": "Error Role required", "status": false,
		})
		return
	}
	if user.Address == "" {
		c.JSON(200, gin.H{
			"message": "Error Address required", "status": false,
		})
		return
	}

	userFound := db.C(UserCollection).Find(bson.M{"email": &user.Email}).One(&user)

	if userFound == nil {
		c.JSON(200, gin.H{
			"message": "User dublicated",
			"status":  false,
		})
		return
	}
	user.Creator = UserInfo.ID
	user.Editor = ""
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.BeforeSave()
	err = db.C(UserCollection).Insert(user)

	if err != nil {
		c.JSON(200, gin.H{
			"message": "Error Insert User",
			"status":  false,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Succes Insert User",
		"status":  true,
	})
}

// Update User Endpoint
func UserInfo(c *gin.Context) {

	authorization := c.Request.Header.Get("Authorization")
	myToken := authorization[7:len(authorization)]

	db := *MongoConfig()
	fmt.Println("MONGO RUNNING: ", db)

	authForm := auth.Auth{}
	userForm := model.User{}

	auth := db.C(AuthCollection).Find(bson.M{"token": &myToken}).One(&authForm)

	if auth != nil {
		c.JSON(200, gin.H{
			"message": "Token not valid",
			"status":  false,
		})
		return
	}

	db.C(UserCollection).Find(bson.M{"_id": &authForm.UserId}).Select(bson.M{"password": 0}).One(&userForm)
	fmt.Println("authForm: ", userForm)

	c.JSON(200, gin.H{
		"message": "Success user found",
		"user":    &userForm,
		"status":  true,
	})
}

func FindUserByToken(c *gin.Context) model.User {

	authorization := c.Request.Header.Get("Authorization")
	myToken := authorization[7:len(authorization)]

	db := *MongoConfig()
	fmt.Println("MONGO RUNNING: ", db)

	authForm := auth.Auth{}
	userForm := model.User{}

	db.C(AuthCollection).Find(bson.M{"token": &myToken}).One(&authForm)
	db.C(UserCollection).Find(bson.M{"_id": &authForm.UserId}).Select(bson.M{"password": 0}).One(&userForm)

	return userForm
}

// Update User Endpoint
func UpdateUser(c *gin.Context) {
	UserInfo := FindUserByToken(c)

	if UserInfo.Role != "admin" {
		if UserInfo.Role != "developer" {
			c.JSON(200, gin.H{
				"message": "You are not authorize",
				"status":  false,
			})
			return
		}
	}
	db := *MongoConfig()
	fmt.Println("MONGO RUNNING: ", db)

	user := model.User{}
	err := c.Bind(&user)

	if err != nil {

		c.JSON(200, gin.H{
			"message": "Error Get Body",
			"status":  false,
		})
		return
	}

	user.Prepare()

	if user.Name == "" {
		c.JSON(200, gin.H{
			"message": "Error Name required", "status": false,
		})
		return
	}
	if user.Email == "" {
		c.JSON(200, gin.H{
			"message": "Error Email required", "status": false,
		})
		return
	}

	if user.Role == "" {
		c.JSON(200, gin.H{
			"message": "Error Role required", "status": false,
		})
		return
	}
	if user.Address == "" {
		c.JSON(200, gin.H{
			"message": "Error Address required", "status": false,
		})
		return
	}

	var id bson.ObjectId = bson.ObjectIdHex(c.Param("id")) // Get Param
	if id == "" {
		c.JSON(200, gin.H{
			"message": "User not valid", "status": false,
		})
		return
	}
	existingUser, err := model.UserInfo(id, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid user"})
		return
	}
	// idParse, errParse := strconv.Atoi(id) // Convert String to Int

	// existingUser.CreatedAt = time.Now()
	existingUser.UpdatedAt = time.Now()
	existingUser.Email = user.Email
	existingUser.Name = user.Name
	existingUser.Address = user.Address
	existingUser.Editor = UserInfo.ID
	if user.Password != "" {
		existingUser.BeforeSave()
	}
	err = db.C(UserCollection).Update(bson.M{"_id": id}, existingUser)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "Error Update User",
			"status":  false,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Succes Update User",
		"user":    &existingUser,
		"status":  true,
	})
}

// Delete User Endpoint
func DeleteUser(c *gin.Context) {
	db := *MongoConfig()
	fmt.Println("MONGO RUNNING: ", db)

	id := c.Param("id")                   // Get Param
	idParse, errParse := strconv.Atoi(id) // Convert String to Int
	if errParse != nil {
		c.JSON(200, gin.H{
			"message": "Error Parse Param",
		})
		return
	}

	err := db.C(UserCollection).Remove(bson.M{"id": &idParse})
	if err != nil {
		c.JSON(200, gin.H{
			"message": "Error Delete User",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Succes Delete User",
	})
}
