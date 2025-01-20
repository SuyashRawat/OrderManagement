package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	CustomerID     string             `json:"cId" bson:"cId"`
	OrdersID       []string           `json:"pid" bson:"_pid"`
	OrdersQuantity []int              `json:"product_quantity" bson:"product_quantity"`
}
