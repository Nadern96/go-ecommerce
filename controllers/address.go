package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadern96/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			c.Abort()
			return
		}
		userObjId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal"})
			c.Abort()
			return
		}

		var address models.Address
		address.ID = primitive.NewObjectID()
		if err = c.BindJSON(&address); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userObjId}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$addresses"}}}}
		group := bson.D{
			{
				Key: "$group",
				Value: bson.D{
					primitive.E{Key: "_id", Value: "$_id"},
					{Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}},
				},
			},
		}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match, unwind, group})
		if err != nil {
			log.Println(err)
		}
		var addressInfo []bson.M
		if err = cursor.All(ctx, &addressInfo); err != nil {
			log.Fatal(err)
		}

		var count int32
		for _, item := range addressInfo {
			count += (item["count"]).(int32)
		}
		if count < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "addresses", Value: address}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				log.Println(err)
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not Allowed"})
		}
		c.JSON(http.StatusOK, "Address Successfully added")
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			c.Abort()
			return
		}

		userObjId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal"})
		}

		var editAddress models.Address
		if err = c.BindJSON(&editAddress); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userObjId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses.0", Value: editAddress}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, "Something went wrong")
			return
		}
		c.JSON(http.StatusOK, "Successfully updated home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			c.Abort()
			return
		}

		userObjId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal"})
		}

		var editAddress models.Address
		if err = c.BindJSON(&editAddress); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userObjId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses.1", Value: editAddress}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, "Something went wrong")
			return
		}
		c.JSON(http.StatusOK, "Successfully updated work address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		userObjId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userObjId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, "Invalid user id")
			return
		}
		ctx.Done()
		c.JSON(200, "Succefully Deleted")
	}
}
