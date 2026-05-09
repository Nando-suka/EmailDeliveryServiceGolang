package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Nando-suka/email-service/model"
	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	client   *redis.Client
	ctx      context.Context
	queueKey string
}

func NewRedisQueue(addr, password, queueName string) *RedisQueue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisQueue{
		client:   rdb,
		ctx:      context.Background(),
		queueKey: queueName,
	}
}

// Enqueue memasukkan task ke antrean (RPUSH)
func (q *RedisQueue) Enqueue(task model.EmailTask) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return q.client.RPush(q.ctx, q.queueKey, data).Err()
}

// Dequeue mengambil task dari antrean (BLPOP, blocking)
func (q *RedisQueue) Dequeue(timeout time.Duration) (*model.EmailTask, error) {
	result, err := q.client.BLPop(q.ctx, timeout, q.queueKey).Result()
	if err != nil {
		return nil, err
	}
	var task model.EmailTask
	err = json.Unmarshal([]byte(result[1]), &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
