package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rubenkristian/backend/commons"
	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/utils"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func InitializeProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

func (ps *ProductService) GetProduct(id uint) (*models.Product, error) {
	var product models.Product

	if err := ps.db.Find(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product with id %d not found", id)
	}

	return &product, nil
}

func (ps *ProductService) GetAllProduct(pagination *commons.PaginationParams) ([]models.Product, error) {
	pagination.SetParams(10, "asc", "id")
	var products []models.Product

	query := ps.db.Model(&models.Product{}).Limit(pagination.Take).Offset(pagination.Skip)

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

func (ps *ProductService) CreateProduct(product *models.Product) error {
	return ps.db.Create(product).Error
}

func (ps *ProductService) UpdateProduct(id uint, input *models.Product) (*models.Product, error) {
	var product models.Product

	if err := ps.db.First(&product, id).Error; err != nil {
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

	if err := ps.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (ps *ProductService) DeleteProduct(id uint) error {
	return ps.db.Delete(&models.Product{}, id).Error
}

func (ps *ProductService) SaveImage(image *multipart.FileHeader) (string, error) {
	if !utils.IsImage(image) {
		return "Bad Request", fmt.Errorf("file is not support, image only")
		// return utils.ResponseError(fiber.StatusBadRequest, "Bad Request", fmt.Errorf("file is not support, image only"))(c)
	}

	os.MkdirAll("./images/product", os.ModePerm)

	randomFileName, err := utils.GenerateImageName(8)

	if err != nil {
		return "Something went wrong", err
		// return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	savePath := filepath.Join("./images/product", fmt.Sprintf("%s-%d.%s", randomFileName, time.Now().Unix(), filepath.Ext(image.Filename)))

	return savePath, nil

	// if err := c.SaveFile(image, savePath); err != nil {
	// 	return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	// }
}
