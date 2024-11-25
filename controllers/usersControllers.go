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
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("users")

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !checkPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Register(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("users")

	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	req.Password = hashedPassword
	req.ID = primitive.NewObjectID().Hex()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	_, err = collection.InsertOne(context.TODO(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": req})
}

func UpdateProfile(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("users")

	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"id": req.ID}
	update := bson.M{
		"$set": bson.M{
			"name":      req.Name,
			"email":     req.Email,
			"address":   req.Address,
			"updatedAt": time.Now(),
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func ChangePassword(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("users")

	var req struct {
		UserID      string `json:"userId"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"id": req.UserID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !checkPasswordHash(req.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	update := bson.M{"$set": bson.M{"password": hashedPassword, "updatedAt": time.Now()}}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"id": req.UserID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func ToggleBlock(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("users")

	userID := c.Param("userId")

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	update := bson.M{"$set": bson.M{"isBlocked": !user.IsBlocked, "updatedAt": time.Now()}}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle block status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User block status updated"})
}

func GetById(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("users")

	userID := c.Param("userId")

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	client := data.GetMongoClient()
	collection := client.Database("foodstoreDB").Collection("users")

	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"id": req.ID}
	update := bson.M{
		"$set": bson.M{
			"name":      req.Name,
			"email":     req.Email,
			"address":   req.Address,
			"isAdmin":   req.IsAdmin,
			"isBlocked": req.IsBlocked,
			"updatedAt": time.Now(),
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
