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

// GetAllFoods retrieves all foods
func GetAllFoods(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("foods")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch foods"})
		return
	}
	defer cursor.Close(context.TODO())

	var foods []models.Food
	if err := cursor.All(context.TODO(), &foods); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse foods"})
		return
	}

	c.JSON(http.StatusOK, foods)
}

// SearchFoods searches for foods by a term
func SearchFoods(c *gin.Context) {
	searchTerm := c.Param("searchTerm")
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("foods")

	filter := bson.M{"name": bson.M{"$regex": searchTerm, "$options": "i"}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch foods"})
		return
	}
	defer cursor.Close(context.TODO())

	var foods []models.Food
	if err := cursor.All(context.TODO(), &foods); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse foods"})
		return
	}

	c.JSON(http.StatusOK, foods)
}

// GetAllTags retrieves all unique tags
func GetAllTags(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("foods")

	cursor, err := collection.Distinct(context.TODO(), "tags", bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, cursor)
}

// GetFoodsByTag retrieves foods by a specific tag
func GetFoodsByTag(c *gin.Context) {
	tag := c.Param("tag")
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("foods")

	filter := bson.M{"tags": tag}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch foods by tag"})
		return
	}
	defer cursor.Close(context.TODO())

	var foods []models.Food
	if err := cursor.All(context.TODO(), &foods); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse foods"})
		return
	}

	c.JSON(http.StatusOK, foods)
}

// GetFoodByID retrieves a food by its ID
func GetFoodByID(c *gin.Context) {
	foodID := c.Param("foodId")
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("foods")

	id, err := primitive.ObjectIDFromHex(foodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	var food models.Food
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&food)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
		return
	}

	c.JSON(http.StatusOK, food)
}

// DeleteFood deletes a food by its ID
func DeleteFood(c *gin.Context) {
	foodID := c.Param("foodId")
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("foods")

	id, err := primitive.ObjectIDFromHex(foodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete food"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food deleted successfully"})
}

// UpdateFood updates an existing food
func UpdateFood(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("foods")

	var food models.Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"_id": food.ID}
	food.UpdatedAt = time.Now()
	update := bson.M{"$set": food}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update food"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food updated successfully"})
}

// AddFood adds a new food item
func AddFood(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("foods")

	var food models.Food

	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	food.ID = primitive.NewObjectID()
	food.CreatedAt = time.Now()
	food.UpdatedAt = time.Now()

	_, err := collection.InsertOne(context.TODO(), food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add food"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food added successfully", "food": food})
}
