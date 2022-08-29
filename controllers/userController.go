package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/joy-adhikary/Restaurant-Management-back_end/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {
		// convert the login data which is on json to somethings that go understand

		// find the user using email and see he exists or not

		// verify the password

		// if all goes well then we generate tokens

		//update tokes and refresh tokens

		//returning it

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		//convert json data coming from web to somethings that go understand

		//validate the data based on user struct

		//check if the email is already been used by another user or not

		// hash the password

		//check if the phone number is already been used by another user or not

		//get some extra details like created_at , updated_at,ID

		//generate token and refresh token (generatealltokens() from helper )

	}
}

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

}
