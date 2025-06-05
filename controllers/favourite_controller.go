package controllers

import (
	"backend_course/common"
	"backend_course/database"
	"backend_course/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type FavouriteController struct {
	Db *gorm.DB
}

func (dbc *FavouriteController) AddToFavourite(c *gin.Context) {
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

	favourite := database.Favourite{
		UserId:    id,
		ProductId: productIdStr,
	}
	if err := dbc.Db.Where("user_id = ? AND product_id = ?", id, productId).FirstOrCreate(&favourite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func (dbc *FavouriteController) DeleteFromFavourite(c *gin.Context) {
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

	if err := dbc.Db.Where("user_id = ? AND product_id = ?", id, productId).Delete(&database.Favourite{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func (dbc *FavouriteController) GetAllFavouriteProducts(c *gin.Context) {
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
		Joins("JOIN favourites ON favourites.product_id = products.id").
		Where("favourites.user_id = ?", id).
		Preload("Images").
		Preload("Carts").
		Preload("Favourites").
		Find(&products).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var result []dto.ProductWithImageResult
	for _, p := range products {
		image := ""
		if len(p.Images) > 0 {
			image = p.Images[0].FilePath
		}

		inCart := false
		for _, cart := range p.Carts {
			if cart.UserId == id {
				inCart = true
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

		result = append(result, dto.ProductWithImageResult{
			Id:         p.Id,
			Name:       p.Name,
			Price:      p.Price,
			CategoryId: p.CategoryId,
			BrandId:    p.BrandId,
			Image:      strings.ReplaceAll(image, "\\", "/"),
			InCart:     inCart,
			IsLiked:    isLiked,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
		"error":  nil,
	})
}
