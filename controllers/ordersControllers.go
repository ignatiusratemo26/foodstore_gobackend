package controllers

import (
	"context"
	"net/http"
	"time"

	"go_backend/data"
	"go_backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateOrderRequest struct {
	Name          string             `json:"name"`
	Address       string             `json:"address"`
	AddressLatLng models.LatLng      `json:"addressLatLng"`
	TotalPrice    float64            `json:"totalPrice"`
	Items         []models.OrderItem `json:"items"`
	UserID        string             `json:"userId"`
}

type PaymentRequest struct {
	PaymentID string `json:"paymentId"`
}

type TrackOrderResponse struct {
	Order  models.Order `json:"order"`
	Status string       `json:"status"`
}

func CreateOrder(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("orders")

	var req models.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = primitive.NewObjectID().Hex()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	req.Status = "Pending"

	_, err := collection.InsertOne(context.TODO(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order": req})
}

func GetNewOrderForCurrentUser(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("orders")

	userID := c.Query("userId")

	var order models.Order
	err := collection.FindOne(context.TODO(), bson.M{"userId": userID, "status": "Pending"}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No new order found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func Pay(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("orders")

	var req struct {
		PaymentID string `json:"paymentId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"paymentId": req.PaymentID}
	update := bson.M{"$set": bson.M{"status": "Paid", "updatedAt": time.Now()}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}

func TrackOrderById(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("orders")

	orderID := c.Param("orderId")

	var order models.Order
	err := collection.FindOne(context.TODO(), bson.M{"id": orderID}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order, "status": order.Status})
}

func GetAll(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("orders")

	state := c.Query("state")

	filter := bson.M{}
	if state != "" {
		filter["status"] = state
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}
	defer cursor.Close(context.TODO())

	var orders []models.Order
	if err = cursor.All(context.TODO(), &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func GetAllStatus(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("orders")

	cursor, err := collection.Distinct(context.TODO(), "status", bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statuses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statuses": cursor})
}
