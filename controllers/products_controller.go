package controllers

import (
	"backend_course/common"
	"backend_course/database"
	"backend_course/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type ProductController struct {
	Db *gorm.DB
}

func (dbc *ProductController) GetAllProducts(c *gin.Context) {
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
	if err := dbc.Db.Preload("Images").Preload("Carts").Preload("Favourites").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  err.Error(),
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

func (dbc *ProductController) GetProductByCategoryId(c *gin.Context) {
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

	categoryId := c.Query("category_id")
	if categoryId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": nil,
			"error":  "invalid category id",
		})
		return
	}

	var products []database.Product
	if err := dbc.Db.Preload("Images").Preload("Carts").Preload("Favourites").Where("category_id = ?", categoryId).
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  err.Error(),
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

func (dbc *ProductController) GetProductInfoById(c *gin.Context) {
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

	productId := c.Query("product_id")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": nil,
			"error":  "product_id is required",
		})
		return
	}

	var product database.Product
	err = dbc.Db.
		Preload("Images").
		Preload("Carts").
		Preload("Favourites").
		Preload("Category").
		Preload("Brand").
		First(&product, productId).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"result": nil, "error": "product not found"})
		return
	}

	var images []string
	for _, img := range product.Images {
		images = append(images, strings.ReplaceAll(img.FilePath, "\\", "/"))
	}

	inCart := false
	for _, cart := range product.Carts {
		if cart.UserId == id {
			inCart = true
			break
		}
	}

	isLiked := false
	for _, fav := range product.Favourites {
		if fav.UserId == id {
			isLiked = true
			break
		}
	}

	result := dto.ProductInfoResult{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category.Name,
		Brand:       product.Brand.Name,
		Images:      images,
		InCart:      inCart,
		IsLiked:     isLiked,
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
		"error":  nil,
	})
}
