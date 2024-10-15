package services

import (
	"context"
	"crm-backend/config"
	"crm-backend/models"
	"encoding/json"
	"errors"
	"strconv"

	"time"

	"github.com/go-redis/redis/v8"
)

func CacheCustomer(customer *models.Customer, redis *redis.Client) error {
	data, err := json.Marshal(customer)
	if err != nil {
		return errors.New("Error marshaling customer data.")
	}

	customerKey := "customer:" + strconv.Itoa(int(customer.ID))
	err = redis.Set(context.Background(), customerKey, data, 5*time.Minute).Err()
	if err != nil {
		return errors.New("Error setting customer data in Redis.")
	}
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
