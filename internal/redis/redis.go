package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis client
type Client struct {
	*redis.Client
}

// NewClient creates a new Redis client
func NewClient(redisURL string) (*Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		// If URL parsing fails, try to connect to localhost
		opts = &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}
	}

	client := redis.NewClient(opts)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("✅ Redis connected successfully")
	return &Client{client}, nil
}

// SetWithTTL sets a key with a TTL
func (c *Client) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value by key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.Client.Exists(ctx, key).Result()
	return result > 0, err
}

// Delete removes a key
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

// Increment increments a counter
func (c *Client) Increment(ctx context.Context, key string) (int64, error) {
	return c.Client.Incr(ctx, key).Result()
}

// SetURL sets a URL in cache with default TTL
func (c *Client) SetURL(ctx context.Context, shortCode, originalURL string) error {
	key := fmt.Sprintf("url:%s", shortCode)
	return c.SetWithTTL(ctx, key, originalURL, 24*time.Hour)
}

// GetURL gets a URL from cache
func (c *Client) GetURL(ctx context.Context, shortCode string) (string, error) {
	key := fmt.Sprintf("url:%s", shortCode)
	return c.Get(ctx, key)
}

// IncrementClickCount increments the click count for a URL
func (c *Client) IncrementClickCount(ctx context.Context, shortCode string) error {
	key := fmt.Sprintf("clicks:%s", shortCode)
	_, err := c.Increment(ctx, key)
	return err
}

// GetClickCount gets the click count for a URL
func (c *Client) GetClickCount(ctx context.Context, shortCode string) (int64, error) {
	key := fmt.Sprintf("clicks:%s", shortCode)
	result, err := c.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	
	// Parse the result as int64
	var count int64
	_, err = fmt.Sscanf(result, "%d", &count)
	return count, err
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.Client.Close()
} 