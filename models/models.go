package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName    *string            `json:"firstName" bson:"firstName" validate:"required,min=2,max=30"`
	LastName     *string            `json:"lastName" bson:"lastName" validate:"required,min=2,max=30"`
	Password     *string            `json:"password" bson:"password" validate:"required,min=6"`
	Email        *string            `json:"email" bson:"email" validate:"email,required"`
	Phone        *string            `json:"phone" bson:"phone" validate:"required"`
	Token        *string            `json:"token" bson:"token"`
	RefreshToken *string            `json:"refreshToken" bson:"refreshToken"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	UserID       string             `json:"userID" bson:"userID"`
	UserCart     []ProductUser      `json:"userCart" bson:"userCart"`
	Addresses    []Address          `json:"addresses" bson:"addresses"`
	Orders       []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Name   *string            `json:"name" bson:"name"`
	Price  *uint64            `json:"price" bson:"price"`
	Rating *uint8             `json:"rating" bson:"rating"`
	Image  *string            `json:"image" bson:"image"`
}

type ProductUser struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Name   *string            `json:"name" bson:"name"`
	Price  *uint64            `json:"price" bson:"price"`
	Rating *uint8             `json:"rating" bson:"rating"`
	Image  *string            `json:"image" bson:"image"`
}

type Address struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	House   *string            `json:"house" bson:"house"`
	Street  *string            `json:"street" bson:"street"`
	City    *string            `json:"city" bson:"city"`
	PinCode *string            `json:"pinCode" bson:"pinCode"`
}

type Order struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Cart          []ProductUser      `json:"cart" bson:"cart"`
	OrderedAt     time.Time          `json:"orderedAt" bson:"orderedAt"`
	Price         int                `json:"price" bson:"price"`
	Discount      *int               `json:"discount" bson:"discount"`
	PaymentMethod Payment            `json:"paymentMethod" bson:"paymentMethod"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"COD" bson:"COD"`
}
