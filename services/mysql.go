package services

import (
	"context"
	"crm-backend/config"
	"crm-backend/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var ctx = context.Background()

func AddCustomer(customer *models.Customer, config *config.Config) error {
	db := config.DB
	if err := db.Create(customer).Error; err != nil {
		return err
	}
	redis := config.Redis
	return CacheCustomer(customer, redis)
}

func GetCustomers(db *gorm.DB) ([]models.Customer, error) {
	var customers []models.Customer
	if err := db.Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func FetchById(id int, db *gorm.DB) (models.Customer, error) {
	var customer models.Customer
	if err := db.First(&customer, id).Error; err != nil {
		return customer, err
	}
	return customer, nil
}

func UpdateCustomer(id int, post models.Customer, config *config.Config) (models.Customer, error) {
	db := config.DB
	customer := models.Customer{
		ID:        post.ID,
		FirstName: post.FirstName,
		LastName:  post.LastName,
		Company:   post.Company,
		City:      post.City,
		Postal:    post.Postal,
		Phone:     post.Phone,
		Email:     post.Email,
		Address:   post.Address,
		Web:       post.Web,
		County:    post.County,
		CreatedAt: post.CreatedAt,
		UpdatedAt: time.Now(),
	}
	if err := db.Model(&models.Customer{}).Where("id = ?", id).Updates(&customer).Error; err != nil {
		return models.Customer{}, err
	}
	redis := config.Redis
	return customer, CacheCustomer(&customer, redis)
}

func DeleteCustomer(id int, config *config.Config) error {
	db := config.DB
	if err := db.Delete(&models.Customer{}, id).Error; err != nil {
		return err
	}

	redisKey := fmt.Sprintf("customer:%d", id)
	redis := config.Redis
	if err := redis.Del(context.Background(), redisKey).Err(); err != nil {
		return err
	}

	return nil
}

func GetAllCustomersFromCache(db *gorm.DB, redisClient *redis.Client) ([]models.Customer, error) {
	var customers []models.Customer

	var customerIDs []int
	if err := db.Model(&models.Customer{}).Pluck("id", &customerIDs).Error; err != nil {
		return nil, err
	}
	for _, id := range customerIDs {
		customerKey := "customer:" + strconv.Itoa(id)
		customerJSON, err := redisClient.Get(ctx, customerKey).Result()
		if err == redis.Nil {
			var customer models.Customer
			if err := db.First(&customer, id).Error; err == nil {
				customerJSON, _ := json.Marshal(customer)
				redisClient.Set(ctx, customerKey, customerJSON, 5*time.Minute)
				customers = append(customers, customer)
			}
		} else if err == nil {
			var customer models.Customer
			if err := json.Unmarshal([]byte(customerJSON), &customer); err == nil {
				customers = append(customers, customer)
			}
		}
	}
	return customers, nil
}
