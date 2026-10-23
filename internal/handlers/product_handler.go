package handlers

import (
	"net/http"
	"strconv"

	"shopProject/internal/database"
	"shopProject/internal/models"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// GetCategories returns all available categories
func (h *ProductHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// CreateCategory adds a category (helper endpoint for populating data)
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"category": category})
}

// GetCategoryProducts returns products under a specific category by ID or by Name
func (h *ProductHandler) GetCategoryProducts(c *gin.Context) {
	identifier := c.Param("identifier")

	var category models.Category

	// Check if identifier is numeric (ID) or string (Category Name)
	categoryID, err := strconv.Atoi(identifier)
	if err == nil {
		if err := database.DB.First(&category, categoryID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
	} else {
		// Search by category name (case-insensitive & supports partial match)
		if err := database.DB.Where("LOWER(name) LIKE LOWER(?)", "%"+identifier+"%").First(&category).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found with name: " + identifier})
			return
		}
	}

	var products []models.Product
	if err := database.DB.Where("category_id = ?", category.ID).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category.Name,
		"products": products,
	})
}

// GetSpecialOffers returns products flagged as special offers
func (h *ProductHandler) GetSpecialOffers(c *gin.Context) {
	var products []models.Product
	if err := database.DB.Where("is_special_offer = ?", true).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch special offers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"special_offers": products})
}

// GetProductDetails returns a single product by ID along with recent comments
func (h *ProductHandler) GetProductDetails(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := database.DB.Preload("Category").Preload("Comments.User").First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// CreateProduct adds a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"product": product})
}

// GetComments returns comments for a specific product
func (h *ProductHandler) GetComments(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var comments []models.Comment
	if err := database.DB.Preload("User").Where("product_id = ?", productID).Order("created_at desc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
}

// AddComment adds a new comment to a product (Authenticated)
func (h *ProductHandler) AddComment(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDVal.(uint)

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify product exists
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	comment := models.Comment{
		ProductID: uint(productID),
		UserID:    userID,
		Content:   req.Content,
		Rating:    req.Rating,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save comment"})
		return
	}

	// Load user info for response
	database.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment submitted successfully",
		"comment": comment,
	})
}
