package handlers

import (
	"net/http"
	"strconv"

	"shopProject/internal/database"
	"shopProject/internal/models"

	"github.com/gin-gonic/gin"
)

type CartHandler struct{}

func NewCartHandler() *CartHandler {
	return &CartHandler{}
}

type AddCartItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// GetCart returns all items in the user's cart
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var cartItems []models.CartItem
	if err := database.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	var totalPrice float64
	for _, item := range cartItems {
		if item.Product != nil {
			price := item.Product.Price
			if item.Product.DiscountPercent > 0 {
				price = price * (1 - (item.Product.DiscountPercent / 100))
			}
			totalPrice += price * float64(item.Quantity)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       cartItems,
		"total_price": totalPrice,
		"count":       len(cartItems),
	})
}

// AddToCart adds a product to the cart or increments quantity if it exists
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify product exists
	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var existingItem models.CartItem
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existingItem)

	if result.RowsAffected > 0 {
		// Update quantity
		existingItem.Quantity += req.Quantity
		if err := database.DB.Save(&existingItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item in cart"})
			return
		}
		database.DB.Preload("Product").First(&existingItem, existingItem.ID)
		c.JSON(http.StatusOK, gin.H{"message": "Cart item updated", "cart_item": existingItem})
		return
	}

	// Create new item
	newItem := models.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := database.DB.Create(&newItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	database.DB.Preload("Product").First(&newItem, newItem.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart", "cart_item": newItem})
}

// RemoveFromCart deletes an item from the cart
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	itemIDStr := c.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", itemID, userID).Delete(&models.CartItem{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in your cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}
