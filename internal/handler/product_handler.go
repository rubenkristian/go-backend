package handler

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/backend/commons"
	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/internal/services"
	"github.com/rubenkristian/backend/utils"
)

type ProductHandler struct {
	productService *services.ProductService
}

func InitializeProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (productHandler *ProductHandler) GetProduct(c *fiber.Ctx) error {
	productId, err := c.ParamsInt("product_id")

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", err)(c)
	}

	product, err := productHandler.productService.GetProduct(uint(productId))

	if err != nil {
		return utils.ResponseError(fiber.StatusNotFound, "Not found", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusOK, "Success get product", product)(c)
}

func (productHandler *ProductHandler) GetAllProduct(c *fiber.Ctx) error {
	paginationParams := &commons.PaginationParams{}

	if err := c.QueryParser(paginationParams); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	products, err := productHandler.productService.GetAllProduct(paginationParams)

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusOK, "Success fetch products", products)(c)
}

func (productHandler *ProductHandler) PostCreateProduct(c *fiber.Ctx) error {
	name := c.FormValue("name")
	desc := c.FormValue("description")
	price, err := strconv.ParseFloat(c.FormValue("price"), 64)

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", err)(c)
	}

	image, err := c.FormFile("image")

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", err)(c)
	}

	productService := productHandler.productService

	result, err := productService.SaveImage(image)

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, result, err)(c)
	}

	if err := c.SaveFile(image, result); err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	var product *models.Product = &models.Product{
		Name:        name,
		Description: desc,
		Price:       price,
		Image:       result,
	}

	productService.CreateProduct(product)

	return utils.ResponseSuccess(fiber.StatusCreated, "Success create product", product)(c)
}

func (productHandler *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	var product models.Product
	imageAvailable := true
	productId, err := c.ParamsInt("product_id")

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	name := c.FormValue("name")
	desc := c.FormValue("description")
	price, err := strconv.ParseFloat(c.FormValue("price"), 64)

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	product.Name = name
	product.Description = desc
	product.Price = price

	image, err := c.FormFile("image")

	if err != nil {
		imageAvailable = false
		product.Image = ""
	}

	if imageAvailable && !utils.IsImage(image) {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	if imageAvailable {
		os.MkdirAll("./images/product", os.ModePerm)

		savePath := filepath.Join("./images/product", image.Filename)

		if err := c.SaveFile(image, savePath); err != nil {
			return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
		}

		product.Image = savePath
	}

	updatedProduct, err := productHandler.productService.UpdateProduct(uint(productId), &product)

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusCreated, "Success update product", updatedProduct)(c)
}

func (productHandler *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	productId, err := c.ParamsInt("product_id")

	if err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", err)(c)
	}

	if err := productHandler.productService.DeleteProduct(uint(productId)); err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusOK, "Success delete product", nil)(c)
}
