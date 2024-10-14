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
)

// func CacheCustomer(customer *models.Customer, redis *redis.Client) error {
// 	// Marshal customer data to JSON format
// 	data, err := json.Marshal(customer)
// 	if err != nil {
// 		fmt.Println("Error marshaling customer data:", err)
// 		return err
// 	}
// 	// Use customer.ID as the Redis key
// 	customerKey := "customer:" + strconv.Itoa(int(customer.ID))

// 	// Set the customer data in Redis with a 5-minute expiration
// 	err = redis.Set(context.Background(), customerKey, data, 5*time.Minute).Err()
// 	if err != nil {
// 		fmt.Println("Error setting customer data in Redis:", err)
// 		return err
// 	}

// 	fmt.Println("Customer cached successfully with key:", customerKey)
// 	return nil
// }

func CacheCustomer(customer *models.Customer, redis *redis.Client) error {
	// Marshal customer data to JSON format
	data, err := json.Marshal(customer)
	if err != nil {
		fmt.Println("Error marshaling customer data:", err)
		return err
	}

	// Use customer.ID as the Redis key
	customerKey := "customer:" + strconv.Itoa(int(customer.ID))

	// Set the customer data in Redis with a 5-minute expiration
	err = redis.Set(context.Background(), customerKey, data, 5*time.Minute).Err()
	if err != nil {
		fmt.Println("Error setting customer data in Redis:", err)
		return err
	}

	fmt.Println("Customer cached successfully with key:", customerKey)
	return nil
}

func GetCachedCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	c := config.Config{}
	val, err := c.Redis.Get(context.Background(), "customers").Result()
	if err == nil {
		json.Unmarshal([]byte(val), &customers)
		return customers, nil
	}
	return customers, err
}
