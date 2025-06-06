package controllers

import (
	"backend_course/common"
	"backend_course/database"
	"backend_course/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type CartController struct {
	Db *gorm.DB
}

func (dbc *CartController) AddToCart(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	token := strings.Split(authHeader, " ")[1]
	claims, err := common.DecodeToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := common.GetIdFromToken(claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	productId := c.Query("product_id")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "product_id is required",
		})
		return
	}

	productIdStr, err := strconv.ParseInt(productId, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	cart := database.Cart{
		UserId:    id,
		ProductId: productIdStr,
		Count:     1,
	}

	if err := dbc.Db.Where("user_id = ? AND product_id = ?", id, productId).FirstOrCreate(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func (dbc *CartController) DeleteFromCart(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	token := strings.Split(authHeader, " ")[1]
	claims, err := common.DecodeToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = common.GetIdFromToken(claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	productId := c.Query("product_id")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "product_id is required",
		})
		return
	}

	if err := dbc.Db.Where("id = ?", productId).Delete(&database.Cart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func (dbc *CartController) GetAllCart(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  "unauthorized",
		})
		return
	}

	token := strings.Split(authHeader, " ")[1]
	claims, err := common.DecodeToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}

	id, err := common.GetIdFromToken(claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  "invalid token",
		})
		return
	}

	var products []database.Product
	if err := dbc.Db.
		Joins("JOIN carts ON carts.product_id = products.id").
		Where("carts.user_id = ?", id).
		Preload("Images").
		Preload("Carts").
		Preload("Favourites").
		Find(&products).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var result []dto.ProductWithCount
	for _, p := range products {
		image := ""
		if len(p.Images) > 0 {
			image = p.Images[0].FilePath
		}

		inCart := false
		count := 0
		var cartId int64
		for _, cart := range p.Carts {
			if cart.UserId == id {
				inCart = true
				count = cart.Count
				cartId = cart.Id
				break
			}
		}

		isLiked := false
		for _, fav := range p.Favourites {
			if fav.UserId == id {
				isLiked = true
				break
			}
		}

		result = append(result, dto.ProductWithCount{
			Id:      cartId,
			Name:    p.Name,
			Price:   p.Price,
			Image:   strings.ReplaceAll(image, "\\", "/"),
			InCart:  inCart,
			IsLiked: isLiked,
			Count:   count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
		"error":  nil,
	})
}

func (dbc *CartController) UpdateProductsCount(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	token := strings.Split(authHeader, " ")[1]
	claims, err := common.DecodeToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = common.GetIdFromToken(claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	cartId := c.Query("cart_id")
	if cartId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cart_id is required",
		})
		return
	}
	count := c.Query("count")
	if count == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "count is required",
		})
		return
	}
	action := c.Query("action")
	if action == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "action is required",
		})
		return
	}

	var cart database.Cart
	if err := dbc.Db.Where("id = ?", cartId).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cart not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch action {
	case "plus":
		cart.Count++
	case "minus":
		cart.Count--
	}

	if err := dbc.Db.Save(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}
