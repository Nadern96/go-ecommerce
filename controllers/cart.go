package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadern96/go-ecommerce/database"
	"github.com/nadern96/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(200, "Successfully added to cart")
	}
}
func (app *Application) RemoveCartItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.userCollection, productID, userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(200, "Successfully remove item from cart")
	}
}
func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.Abort()
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid user id"})
			return
		}

		userObjId, _ := primitive.ObjectIDFromHex(userId)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userObjId}}).Decode(&filledCart)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, "Not found")
			return
		}

		match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userObjId}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$userCart"}}}}
		group := bson.D{
			{
				Key: "$group",
				Value: bson.D{
					primitive.E{Key: "_id", Value: "$_id"},
					{Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$userCart.price"}}},
				},
			},
		}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match, unwind, group})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = cursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		var total int64
		for _, item := range listing {
			total += (item["total"]).(int64)
		}

		res := struct {
			Total int64
			Cart  []models.ProductUser
		}{
			Total: total,
			Cart:  filledCart.UserCart,
		}

		c.JSON(200, res)
	}
}
func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.Panic("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(200, "Successfully placed the order")
	}
}
func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuy(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(200, "Successfully placed the order")
	}
}
