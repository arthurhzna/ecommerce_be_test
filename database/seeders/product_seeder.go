package seeders

import (
	"github.com/arthurhzna/ecommerce_be_test/domain/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RunProductSeeder(db *gorm.DB) {
	product := models.Product{
		UUID:        uuid.New(),
		Name:        "Laptop Lenovo",
		Description: "Laptop Lenovo ThinkPad X1 Carbon",
		Price:       10000000,
	}

	err := db.FirstOrCreate(&product, models.Product{Name: product.Name}).Error
	if err != nil {
		logrus.Errorf("failed to seed product: %v", err)
		panic(err)
	}
	logrus.Infof("product %s successfully seeded", product.Name)
}
