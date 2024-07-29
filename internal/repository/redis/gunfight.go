package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"math"
	"strconv"
)

type GunfightRedisRepository struct {
	BaseRedis
}

func NewGunfightRedis(redis *redis.Client) *GunfightRedisRepository {
	return &GunfightRedisRepository{
		BaseRedis: BaseRedis{redis: redis},
	}
}

// AddPlayerToQueue добавляет игрока в очередь для поиска соперника
func (r *GunfightRedisRepository) AddPlayerToQueue(ctx context.Context, userID int, gold int) error {
	return r.redis.ZAdd(ctx, "gunfight_queue", redis.Z{
		Score:  float64(gold),
		Member: userID,
	}).Err()
}

// FindOpponent ищет соперника для игрока с заданным количеством золота
func (r *GunfightRedisRepository) FindOpponent(ctx context.Context, gold int) (int, error) {
	minGold := int(math.Round(float64(gold) * 0.95))
	maxGold := int(math.Round(float64(gold) * 1.05))

	opponents, err := r.redis.ZRangeByScore(ctx, "gunfight_queue", &redis.ZRangeBy{
		Min:    strconv.Itoa(minGold), // Преобразование целого числа в строку
		Max:    strconv.Itoa(maxGold), // Преобразование целого числа в строку
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		return 0, err
	}

	if len(opponents) == 0 {
		return 0, nil // No opponent found
	}

	opponentID, err := strconv.Atoi(opponents[0])
	if err != nil {
		return 0, err
	}

	return opponentID, nil
}

// RemovePlayerFromQueue Удаляет игрока из очереди
func (r *GunfightRedisRepository) RemovePlayerFromQueue(ctx context.Context, userID int) error {
	return r.redis.ZRem(ctx, "gunfight_queue", userID).Err()
}
