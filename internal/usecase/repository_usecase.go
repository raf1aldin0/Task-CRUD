package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Task-CRUD/internal/cbreaker"
	"Task-CRUD/internal/entity"
	interfaces "Task-CRUD/internal/interfaces"

	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/sony/gobreaker"
)

type RepoUseCase struct {
	repoRepo interfaces.RepoRepositoryInterfaceGorm
	redis    *redis.Client
	breaker  *gobreaker.CircuitBreaker
	kafka    *kafka.Writer
	validate *validator.Validate
}

func NewRepoUseCaseFull(
	repoRepo interfaces.RepoRepositoryInterfaceGorm,
	redisClient *redis.Client,
	kafkaWriter *kafka.Writer,
) interfaces.RepoUseCaseInterface {
	return &RepoUseCase{
		repoRepo: repoRepo,
		redis:    redisClient,
		breaker:  cbreaker.Breaker,
		kafka:    kafkaWriter,
		validate: validator.New(),
	}
}

func (uc *RepoUseCase) GetAllRepos(ctx context.Context) ([]entity.Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.GetAllRepos")
	defer span.Finish()

	cacheKey := "repositories:all"

	if uc.redis != nil {
		if cached, err := uc.redis.Get(ctx, cacheKey).Result(); err == nil {
			var repos []entity.Repository
			if err := json.Unmarshal([]byte(cached), &repos); err == nil {
				span.LogFields(log.String("cache", "hit"))
				fmt.Println("✅ Data repositories diambil dari Redis")
				return repos, nil
			}
		}
	}

	result, err := uc.breaker.Execute(func() (interface{}, error) {
		return uc.repoRepo.GetAllRepositories(ctx)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, fmt.Errorf("get all repositories failed: %w", err)
	}

	repos := result.([]entity.Repository)

	if uc.redis != nil {
		bytes, _ := json.Marshal(repos)
		_ = uc.redis.Set(ctx, cacheKey, bytes, 10*time.Minute).Err()
	}

	return repos, nil
}

func (uc *RepoUseCase) GetRepositoryByID(ctx context.Context, id uint) (*entity.Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.GetRepositoryByID")
	defer span.Finish()

	cacheKey := fmt.Sprintf("repository:%d", id)

	if uc.redis != nil {
		if cached, err := uc.redis.Get(ctx, cacheKey).Result(); err == nil {
			var repo entity.Repository
			if err := json.Unmarshal([]byte(cached), &repo); err == nil {
				span.LogFields(log.String("cache", "hit"))
				fmt.Println("✅ Repository ditemukan di Redis")
				return &repo, nil
			}
		}
	}

	result, err := uc.breaker.Execute(func() (interface{}, error) {
		return uc.repoRepo.GetRepositoryByID(ctx, id)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, fmt.Errorf("get repository by ID failed: %w", err)
	}

	repo := result.(*entity.Repository)

	if uc.redis != nil {
		bytes, _ := json.Marshal(repo)
		_ = uc.redis.Set(ctx, cacheKey, bytes, 10*time.Minute).Err()
	}

	return repo, nil
}

func (uc *RepoUseCase) CreateRepo(ctx context.Context, repo *entity.Repository) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.CreateRepo")
	defer span.Finish()

	if err := uc.validate.Struct(repo); err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := uc.breaker.Execute(func() (interface{}, error) {
		return nil, uc.repoRepo.CreateRepository(ctx, repo)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("create repository failed: %w", err)
	}

	if uc.redis != nil {
		_ = uc.redis.Del(ctx, "repositories:all").Err()
	}

	return uc.sendKafkaMessage(ctx, "repository_created", repo)
}

func (uc *RepoUseCase) UpdateRepo(ctx context.Context, id uint, repo *entity.Repository) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.UpdateRepo")
	defer span.Finish()

	if err := uc.validate.Struct(repo); err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := uc.breaker.Execute(func() (interface{}, error) {
		return nil, uc.repoRepo.UpdateRepository(ctx, id, repo)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("update repository failed: %w", err)
	}

	if uc.redis != nil {
		pipeline := uc.redis.TxPipeline()
		pipeline.Del(ctx, "repositories:all")
		pipeline.Del(ctx, fmt.Sprintf("repository:%d", id))
		_, _ = pipeline.Exec(ctx)
	}

	return uc.sendKafkaMessage(ctx, "repository_updated", repo)
}

func (uc *RepoUseCase) DeleteRepo(ctx context.Context, id uint) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepoUseCase.DeleteRepo")
	defer span.Finish()

	_, err := uc.breaker.Execute(func() (interface{}, error) {
		return nil, uc.repoRepo.DeleteRepository(ctx, id)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return fmt.Errorf("delete repository failed: %w", err)
	}

	if uc.redis != nil {
		pipeline := uc.redis.TxPipeline()
		pipeline.Del(ctx, "repositories:all")
		pipeline.Del(ctx, fmt.Sprintf("repository:%d", id))
		_, _ = pipeline.Exec(ctx)
	}

	return uc.sendKafkaMessage(ctx, "repository_deleted", map[string]uint{"id": id})
}

func (uc *RepoUseCase) sendKafkaMessage(ctx context.Context, topic string, payload interface{}) error {
	if uc.kafka == nil {
		return nil
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal kafka payload failed: %w", err)
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(topic),
		Value: bytes,
	}
	if err := uc.kafka.WriteMessages(ctx, msg); err != nil {
		fmt.Println("❌ Kafka send failed:", err)
		return err
	}
	fmt.Println("📤 Kafka event terkirim:", topic)
	return nil
}
