package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/DiasOrazbaev/url-shortner/shortner"
	"github.com/go-redis/redis/v9"
	"strconv"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisUrl string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedisRepository(redisUrl string) (shortner.RedirectRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisUrl)
	if err != nil {
		return nil, errors.New(err.Error() + " repository.redis.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*shortner.Redirect, error) {
	redirect := &shortner.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, errors.New(err.Error() + " repository.redis.redisRepository.Find")
	}
	if len(data) == 0 {
		return nil, errors.New(shortner.ErrRedirectNotFound.Error() + " repository.redis.redisRepository.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.New(err.Error() + " repository.redis.redisRepository.Find")
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r *redisRepository) Save(redirect *shortner.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	_, err := r.client.HSet(context.Background(), key, data).Result()
	if err != nil {
		return errors.New(err.Error() + " repository.redis.redisRepository.Save")
	}

	return nil
}
