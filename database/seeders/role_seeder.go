package seeders

import (
	"github.com/arthurhzna/ecommerce_be_test/constants"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RunRoleSeeder(db *gorm.DB) {
	roles := []models.Role{
		{
			ID:   constants.Admin,
			Code: "ADMIN",
			Name: "Administrator",
		},
		{
			ID:   constants.Customer,
			Code: "CUSTOMER",
			Name: "Customer",
		},
	}

	for _, role := range roles {
		var existingRole models.Role
		result := db.First(&existingRole, role.ID)

		if result.Error == nil {
			if existingRole.Code != role.Code || existingRole.Name != role.Name {
				existingRole.Code = role.Code
				existingRole.Name = role.Name
				if err := db.Save(&existingRole).Error; err != nil {
					logrus.Errorf("failed to update role: %v", err)
					panic(err)
				}
			}
			logrus.Infof("role %s already exists", role.Code)
		} else {
			if err := db.Create(&role).Error; err != nil {
				logrus.Errorf("failed to seed role: %v", err)
				panic(err)
			}
			logrus.Infof("role %s successfully seeded", role.Code)
		}
	}
}
