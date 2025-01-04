package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/backend/commons"
	"github.com/rubenkristian/backend/configs"
	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/utils"
	"gorm.io/gorm"
)

type ProductService struct {
	db        *gorm.DB
	s3Service *utils.S3Service
	s3Config  *configs.S3Config
}

func InitializeProductService(db *gorm.DB, s3Service *utils.S3Service, s3Config *configs.S3Config) *ProductService {
	return &ProductService{
		db:        db,
		s3Service: s3Service,
		s3Config:  s3Config,
	}
}

func (productService *ProductService) GetProduct(id uint) (*models.Product, error) {
	var product models.Product

	if err := productService.db.Find(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product with id %d not found", id)
	}

	return &product, nil
}

func (productService *ProductService) GetAllProduct(pagination *commons.PaginationParams) ([]models.Product, error) {
	pagination.SetParams(10, "asc", "id")
	var products []models.Product

	query := productService.db.Model(&models.Product{}).Limit(pagination.Take).Offset(pagination.Skip)

	trimSearch := strings.TrimSpace(pagination.Search)

	if trimSearch != "" {
		query = query.Where("name LIKE ?", "%"+trimSearch+"%")
	}

	err := query.Order(pagination.SortBy + " " + pagination.Sort).Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (productService *ProductService) CreateProduct(product *models.Product) error {
	return productService.db.Create(product).Error
}

func (productService *ProductService) UpdateProduct(id uint, input *models.Product) (*models.Product, error) {
	var product models.Product

	if err := productService.db.First(&product, id).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if input.Name != "" {
		product.Name = input.Name
	}

	if input.Description != "" {
		product.Description = input.Description
	}

	if input.Price != 0 {
		product.Price = input.Price
	}

	if err := productService.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (productService *ProductService) DeleteProduct(id uint) error {
	return productService.db.Delete(&models.Product{}, id).Error
}

func (productService *ProductService) SaveImage(c *fiber.Ctx, file string) (string, error) {
	image, err := c.FormFile("image")

	if err != nil {
		return "Something went wrong", err
	}

	fileType, ok := utils.IsImage(image)

	if !ok {
		return "Bad Request", fmt.Errorf("file is not support, image only")
	}

	os.MkdirAll("./images/product", os.ModePerm)

	randomFileName, err := utils.GenerateImageName(8)

	if err != nil {
		return "Something went wrong", err
	}

	savePath := filepath.Join("./images/product", fmt.Sprintf("%s-%d.%s", randomFileName, time.Now().Unix(), filepath.Ext(image.Filename)))

	if err := c.SaveFile(image, savePath); err != nil {
		return "Something went wrong", err
	}

	result := make(chan error)

	go productService.s3Service.UploadFile(savePath, fileType, result)

	err = <-result

	if err != nil {
		return "Something went wrong", err
	}

	return randomFileName, nil
}
