package model

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/guizot/go-gin-mongodb/config"
	"github.com/guizot/go-gin-mongodb/src/security"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//User
type User struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string        `bson:"email"`
	Name      string        `bson:"name"`
	Address   string        `bson:"address"`
	Password  string        `bson:"password"`
	Role      string        `bson:"role"`
	Age       int           `bson:"age"`
	Creator   bson.ObjectId `json:"creator,omitempty" bson:"creator,omitempty"`
	Editor    bson.ObjectId `json:"editor,omitempty" bson:"editor,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
func MongoConfig() *mgo.Database {
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
	}
	return db
}
func UserInfo(id bson.ObjectId, userCollection string) (User, error) {
	// Get DB from Mongo Config
	db := *MongoConfig()
	user := User{}
	err := db.C(userCollection).Find(bson.M{"_id": &id}).One(&user)
	return user, err
}

func (u *User) Prepare() {

	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Address = html.EscapeString(strings.TrimSpace(u.Address))
}

type Users []User
