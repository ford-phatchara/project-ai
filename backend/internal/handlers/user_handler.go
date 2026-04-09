package handlers

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/database"
// 	"github.com/phatcharasangsuphap/gemlni-cli-backend/internal/models"
// )

// GetUsers handles GET /users request
// func GetUsers(c *gin.Context) {
// 	var users []models.User
// 	database.DB.Find(&users)

// 	c.JSON(http.StatusOK, gin.H{"data": users})
// }

// GetUserById handles GET /users/:id request
// func GetUserById(c *gin.Context) {
// 	var user models.User
// 	if err := database.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": user})
// }

// CreateUser handles POST /users request
// func CreateUser(c *gin.Context) {
// 	var input models.User
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	user := models.User{Name: input.Name, Email: input.Email}
// 	database.DB.Create(&user)

// 	c.JSON(http.StatusOK, gin.H{"data": user})
// }

// UpdateUser handles PUT /users/:id request
// func UpdateUser(c *gin.Context) {
// 	var user models.User
// 	if err := database.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
// 		return
// 	}

// 	var input struct {
// 		Name  string `json:"name"`
// 		Email string `json:"email"`
// 	}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Use Map or explicit fields to avoid GORM issues with zero values or unintended updates
// 	database.DB.Model(&user).Updates(models.User{Name: input.Name, Email: input.Email})

// 	c.JSON(http.StatusOK, gin.H{"data": user})
// }
