package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age int    `json:"age"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Age: 25},
	{ID: 2, Name: "Jane Doe", Age: 30},
}

func main() {
	// r := gin.Default()
	// r.GET("/users", GetUsers)
	// r.GET("/users/:id", GetUser)
	// r.POST("/users", CreateUser)
	// r.PUT("/users/:id", UpdateUser)
	// r.DELETE("/users/:id", DeleteUser)
	// r.Run("localhost:5000")
	
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	users = append(users, user)
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	for i, u := range users {
		if u.ID == id {
			users[i] = user
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusNoContent, nil)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

