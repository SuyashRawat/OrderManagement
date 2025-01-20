package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"order/models"
	"order/producer"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository interface defines the methods we need for DB operations
type OrderRepository struct {
	client *mongo.Client
}

// NewUserRepository creates and returns a new UserRepository instance
func NewOrderRepository(client *mongo.Client) *OrderRepository {
	return &OrderRepository{client}
}

var db string = "service1"
var c string = "orders"

func (ur *OrderRepository) CreateOrder(order *models.Order) error {
	collection := ur.client.Database(db).Collection(c)

	if order.ID.IsZero() {
		order.ID = primitive.NewObjectID()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//validate order
	ind := 0
	id := order.OrdersID[ind]
	quan := order.OrdersQuantity[ind]
	conn, channel, err := producer.MQConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer channel.Close()

	// Create a sample registration data
	regData := producer.Order{
		ProductId: id,
		Quantity:  quan,
	}

	// Marshal the data into JSON
	message, err := json.Marshal(regData)
	if err != nil {
		log.Fatal(err)
	}

	// Send the message to RabbitMQ
	err = producer.MQPublish(channel, message)
	if err != nil {
		log.Fatal(err)
	}

	// Log successful message sending
	fmt.Println("Message sent:", string(message))
	_, err = collection.InsertOne(ctx, order)
	return err
}
func (ur *OrderRepository) CreateBulkOrder(order *models.Order) error {
	collection := ur.client.Database(db).Collection(c)

	if order.ID.IsZero() {
		order.ID = primitive.NewObjectID()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, channel, err := producer.MQConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer channel.Close()
	//validate order
	for ind, _ := range order.OrdersID {
		go func() {
			//pass values to rabbitmq queue
			id := order.OrdersID[ind]
			quan := order.OrdersQuantity[ind]

			// Create a sample registration data
			regData := producer.Order{
				ProductId: id,
				Quantity:  quan,
				UserId:    order.CustomerID,
			}

			// Marshal the data into JSON
			message, err := json.Marshal(regData)
			if err != nil {
				log.Fatal(err)
			}

			// Send the message to RabbitMQ
			err = producer.MQPublish(channel, message)
			if err != nil {
				log.Fatal(err)
			}

			// Log successful message sending
			fmt.Println("Message sent:", string(message))
		}()
	}
	_, err = collection.InsertOne(ctx, order)
	return err
}
func (ur *OrderRepository) GetOrder(id primitive.ObjectID) (*models.Order, error) {
	collection := ur.client.Database(db).Collection(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var order models.Order
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetUsers fetches all users from the database
func (ur *OrderRepository) GetOrders() ([]models.Order, error) {
	collection := ur.client.Database(db).Collection(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	for cursor.Next(ctx) {
		var user models.Order
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		orders = append(orders, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// UpdateUser updates a user's details by ID
func (ur *OrderRepository) UpdateOrder(id primitive.ObjectID, user *models.Order) error {
	collection := ur.client.Database(db).Collection(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = id
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": user})
	return err
}

// DeleteUser deletes a user by ID
func (ur *OrderRepository) DeleteOrder(id primitive.ObjectID) error {
	collection := ur.client.Database(db).Collection(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
