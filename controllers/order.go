package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"order/models"
	"order/repository"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderController struct {
	orderRepo *repository.OrderRepository
}

func NewOrderController(userRepo *repository.OrderRepository) *OrderController {
	return &OrderController{userRepo}
}

// GetUser retrieves a user by ID
func (uc *OrderController) GetOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	order, err := uc.orderRepo.GetOrder(objectID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(uj)
}

// GetUsers retrieves all users
func (uc *OrderController) GetOrders(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	orders, err := uc.orderRepo.GetOrders()
	if err != nil {
		http.Error(w, "Error finding orders", http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(uj)
}

// CreateUser creates a new user
func (uc *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var order models.Order

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&order)

	if err != nil {
		http.Error(w, "request body invalid", http.StatusBadRequest)
		return
	}
	if order.CustomerID == "" || order.OrdersID == nil || order.OrdersQuantity == nil {

		http.Error(w, "request body invalid(some of the fields of the body are not provided)", http.StatusBadRequest)
		return
	}

	if len(order.OrdersID) != 1 || len(order.OrdersQuantity) != 1 {
		http.Error(w, "request body invalid(only 1 pid and 1 product quantity needed)", http.StatusBadRequest)
		return
	}
	err = uc.orderRepo.CreateBulkOrder(&order)
	if err != nil {
		msg := "error inserting order with orderid :" + err.Error()
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	uj, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(uj)
}
func (uc *OrderController) CreateBulkOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var order models.Order

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&order)

	if err != nil {
		http.Error(w, "request body invalid", http.StatusBadRequest)
		return
	}
	if order.CustomerID == "" || order.OrdersID == nil || order.OrdersQuantity == nil {

		http.Error(w, "request body invalid(some of the fields of the body are not provided)", http.StatusBadRequest)
		return
	}

	if len(order.OrdersID) != len(order.OrdersQuantity) {
		http.Error(w, "length of order id and order quantity not the same", http.StatusBadRequest)
		return
	}

	err = uc.orderRepo.CreateBulkOrder(&order)
	if err != nil {
		msg := "error inserting order with orderid :" + err.Error()
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	uj, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(uj)
}

// UpdateUser updates a user's details
func (uc *OrderController) UpdateOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "invalid order id for updation", http.StatusBadRequest)
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = uc.orderRepo.UpdateOrder(objectID, &order)
	if err != nil {
		http.Error(w, "Error updating order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated order with ID: %s\n", id)
}

// DeleteUser deletes a user
func (uc *OrderController) DeleteOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	_, err = uc.orderRepo.GetOrder(objectID)
	if err != nil {
		http.Error(w, "no order with this ID", http.StatusNotFound)
		return
	}

	err = uc.orderRepo.DeleteOrder(objectID)
	if err != nil {
		http.Error(w, "Error deleting order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted order with ID: %s\n", id)
}
