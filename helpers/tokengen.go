package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nadern96/go-ecommerce/database"
	"github.com/nadern96/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")
var UserCollection *mongo.Collection = database.UserData(database.Client, "users")

func TokenGenerator(user *models.User) (token, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:     *user.Email,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Uid:       *&user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
func ValidateToken(clientToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		clientToken,
		&SignedDetails{},
		func(clientToken *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}

func UpdateAllTokens(token, refreshToken, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: token})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: refreshToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format((time.RFC3339)))

	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: updatedAt})
	upsert := true
	filter := bson.M{"userID": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := UserCollection.UpdateOne(
		ctx,
		filter,
		bson.D{{Key: "$set", Value: updateObj}},
		&opt,
	)
	if err != nil {
		log.Panic(err)
		return
	}
	return
}
