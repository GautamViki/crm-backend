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
	// Cache the data in Redis
	return CacheCustomer(customer, redis)
}

func GetCustomers(db *gorm.DB) ([]models.Customer, error) {
	var customers []models.Customer
	if err := db.Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func UpdateCustomer(id int, customer *models.Customer, config *config.Config) error {
	db := config.DB

	// Update the customer in MySQL
	if err := db.Model(&models.Customer{}).Where("id = ?", id).Updates(customer).Error; err != nil {
		return err
	}

	// Update the Redis cache with the updated customer data
	redis := config.Redis
	return CacheCustomer(customer, redis) // Update Redis cache
}

func DeleteCustomer(id int, config *config.Config) error {
	db := config.DB
	if err := db.Delete(&models.Customer{}, id).Error; err != nil {
		fmt.Println("Error deleting customer from database:", err)
		return err
	}

	// Step 2: Remove customer data from Redis cache
	redisKey := fmt.Sprintf("customer:%d", id)
	redis := config.Redis
	if err := redis.Del(context.Background(), redisKey).Err(); err != nil {
		fmt.Println("Error deleting customer from Redis cache:", err)
		return err
	}

	fmt.Printf("Customer with ID %d deleted successfully.\n", id)
	return nil
}

func GetAllCustomersFromCache(db *gorm.DB, redisClient *redis.Client) ([]models.Customer, error) {
	var customers []models.Customer

	// Retrieve all customer IDs from MySQL
	var customerIDs []int
	if err := db.Model(&models.Customer{}).Pluck("id", &customerIDs).Error; err != nil {
		return nil, err // Return error if not found in MySQL
	}

	// Check each customer in Redis
	for _, id := range customerIDs {
		customerKey := "customer:" + strconv.Itoa(id)
		customerJSON, err := redisClient.Get(ctx, customerKey).Result()
		if err == redis.Nil {
			// If customer not found in Redis, fetch from MySQL
			var customer models.Customer
			if err := db.First(&customer, id).Error; err == nil {
				// Cache the customer in Redis
				customerJSON, _ := json.Marshal(customer)
				redisClient.Set(ctx, customerKey, customerJSON, 5*time.Minute)
				customers = append(customers, customer) // Add to the results
			}
		} else if err == nil {
			// If customer found in Redis, unmarshal and add to results
			var customer models.Customer
			if err := json.Unmarshal([]byte(customerJSON), &customer); err == nil {
				customers = append(customers, customer)
			}
		}
	}

	return customers, nil
}
