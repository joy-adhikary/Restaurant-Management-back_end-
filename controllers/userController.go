package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joy-adhikary/Restaurant-Management-back_end/database"
	"github.com/joy-adhikary/Restaurant-Management-back_end/helpers"
	"github.com/joy-adhikary/Restaurant-Management-back_end/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		// here we need to set the checker if we want to fix the page reload

		result, err := userCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "data fatch issue arise "})
			return

		}

		var allusers []bson.M

		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allusers)

	}
}

func GetUser() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occure when fatch one data from get user "})
		}

		c.JSON(http.StatusOK, user)

	}
}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var user models.User

		// convert the login data which is on json to somethings that go understand
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// find the user using email and see he exists or not
		var foundUser models.User

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user email not found !! login again "})
			return
		}

		// verify the password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancle()

		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// if all goes well then we generate tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id) // token banabe email,fname,lname r userid er upr base kore

		//update tokes and refresh tokens

		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		//status return

		c.JSON(http.StatusOK, foundUser)

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancle = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancle()

		var user models.User

		//convert json data coming from web to somethings that go understand
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//validate the data based on user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
			return
		}

		//check if the email is already been used by another user or not

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured when email check"})
			return
		}

		// hash the password

		password := HashPassword(*user.Password)
		user.Password = &password

		//check if the phone number is already been used by another user or not

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured when phone check"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user already exsits"})
			return
		}

		//get some extra details like created_at , updated_at,ID

		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		//generate token and refresh token (generatealltokens() from helper )
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id) // token banabe email,fname,lname r userid er upr base kore
		user.Token = &token
		user.Refresh_Token = &refreshToken

		resultInsertionNumber, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user item was not created"})
			return
		}
		defer cancle()

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func HashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	msg := " "
	check := true
	if err != nil {
		msg = fmt.Sprintf("password missmatched")
		check = false
	}
	return check, msg
}
