package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nadern96/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProduct  = errors.New("can't find the product")
	ErrUserIdIsNotValid   = errors.New("invalid user")
	ErrCantUpdateUser     = errors.New("can't add product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove item from the cart")
	ErrCantGetItem        = errors.New("can't get item from the cart")
	ErrCantBuyCartItem    = errors.New("can't update the purchase")
)

func AddToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userId string) error {
	searchFromDb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchFromDb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProduct
	}
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: userObjId}}
	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				primitive.E{
					Key: "userCart",
					Value: bson.D{
						{
							Key:   "$each",
							Value: productCart,
						},
					},
				},
			},
		},
	}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}
func RemoveCartItem(ctx context.Context, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	userObjId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: userObjId}}
	update := bson.M{"$pull": bson.M{"userCart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil
}

// func GetItemFromCart() error {

// }
func BuyFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	userObjId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var user models.User
	var order models.Order

	order.ID = primitive.NewObjectID()
	order.OrderedAt = time.Now()
	order.Cart = make([]models.ProductUser, 0)
	order.PaymentMethod.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$userCart"}}}}
	group := bson.D{
		{
			Key: "$group",
			Value: bson.D{
				primitive.E{Key: "_id", Value: "$_id"},
				{
					Key:   "total",
					Value: bson.D{primitive.E{Key: "$sum", Value: "$userCart.price"}},
				},
			},
		},
	}
	res, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, group})
	if err != nil {
		log.Fatal(err)
	}

	var userCart []bson.M
	if err = res.All(ctx, &userCart); err != nil {
		log.Fatal(err)
	}

	var totalPrice int64
	for _, item := range userCart {
		price := item["total"]
		totalPrice += price.(int64)
	}
	order.Price = int(totalPrice)
	filter := bson.D{primitive.E{Key: "_id", Value: userObjId}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userObjId}}).Decode(&user)
	if err != nil {
		log.Println(err)
	}

	filter = bson.D{primitive.E{Key: "_id", Value: userObjId}}
	update = bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders.$[].order_list", Value: bson.M{"$each": user.UserCart}}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	userCartEmpty := make([]models.ProductUser, 0)
	filter = bson.D{primitive.E{Key: "_id", Value: userObjId}}
	update = bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "userCart", Value: userCartEmpty}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuy(ctx context.Context, prodCollection, userCollection *mongo.Collection, prodID primitive.ObjectID, userID string) error {
	userObjId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var prodDetails models.ProductUser
	var orderDetails models.Order

	orderDetails.ID = primitive.NewObjectID()
	orderDetails.OrderedAt = time.Now()
	orderDetails.Cart = make([]models.ProductUser, 0)
	orderDetails.PaymentMethod.COD = true

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: prodID}}).Decode(&prodDetails)
	if err != nil {
		log.Println(err)
	}

	orderDetails.Price = int(*prodDetails.Price)
	filter := bson.D{primitive.E{Key: "_id", Value: userObjId}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderDetails}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter = bson.D{primitive.E{Key: "_id", Value: userObjId}}
	update = bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders.$[].order_list", Value: orderDetails}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	return nil
}
