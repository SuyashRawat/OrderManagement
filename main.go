package main

import (
	"log"
	"net/http"

	"order/config"
	"order/controllers"
	"order/repository"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Connect to MongoDB
	client, err := config.ConnectToDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// fmt.Print(client)

	// Create a new router
	router := httprouter.New()

	// Initialize the controller
	orderRepo := repository.NewOrderRepository(client)
	orderController := controllers.NewOrderController(orderRepo)

	// Define the routes
	router.GET("/order/:id", orderController.GetOrder)
	router.GET("/order", orderController.GetOrders)
	router.POST("/order", orderController.CreateOrder)
	router.POST("/bulkorder", orderController.CreateBulkOrder)
	router.PUT("/order/:id", orderController.UpdateOrder)
	router.DELETE("/order/:id", orderController.DeleteOrder)

	// Start the server
	log.Println("Starting server on :9001")
	log.Fatal(http.ListenAndServe(":9001", router))
	// Connect to RabbitMQ
	// conn, channel, err := MQConnect()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()
	// defer channel.Close()

	// // Create a sample registration data
	// regData := Order{
	// 	ProductId: "user@example.com",
	// 	Quantity:  1,
	// }

	// // Marshal the data into JSON
	// message, err := json.Marshal(regData)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Send the message to RabbitMQ
	// err = MQPublish(channel, message)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Log successful message sending
	// fmt.Println("Message sent:", string(message))
}
