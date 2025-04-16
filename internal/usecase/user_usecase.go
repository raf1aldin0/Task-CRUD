package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"Task-CRUD/internal/cbreaker"
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/interfaces"

	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

type UserUseCase struct {
	userRepo interfaces.UserRepositoryInterfaceSQL  // Menggunakan UserRepositoryInterfaceSQL
	redis    *redis.Client
	breaker  *gobreaker.CircuitBreaker
	validate *validator.Validate
}

func NewUserUseCase(userRepo interfaces.UserRepositoryInterfaceSQL, validate *validator.Validate) interfaces.UserUseCaseInterface {
	return &UserUseCase{
		userRepo: userRepo,
		breaker:  cbreaker.Breaker,
		validate: validate,
	}
}

func NewUserUseCaseWithCache(userRepo interfaces.UserRepositoryInterfaceSQL, redisClient *redis.Client, validate *validator.Validate) interfaces.UserUseCaseInterface {
	return &UserUseCase{
		userRepo: userRepo,
		redis:    redisClient,
		breaker:  cbreaker.Breaker,
		validate: validate,
	}
}

func (uc *UserUseCase) GetUsers(ctx context.Context) ([]entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.GetUsers")
	defer span.Finish()

	span.SetTag("operation", "GetUsers")
	span.SetTag("source", "UserService")
	span.SetTag("cache", "miss")

	cacheKey := "users:all"

	if uc.redis != nil {
		cached, err := uc.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var users []entity.User
			if err := json.Unmarshal([]byte(cached), &users); err == nil {
				span.LogFields(log.String("cache", "hit"))
				span.SetTag("cache", "hit")
				fmt.Println("✅ Data users diambil dari Redis")
				return users, nil
			}
			span.LogFields(log.Error(err))
			fmt.Printf("⚠️ Gagal unmarshal data users dari Redis: %v\n", err)
		} else if err != redis.Nil {
			span.LogFields(log.Error(err))
			fmt.Printf("⚠️ Redis error: %v\n", err)
		}
	}

	result, err := uc.breaker.Execute(func() (interface{}, error) {
		return uc.userRepo.GetAllUsers(ctx)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, err
	}
	users := result.([]entity.User)

	if uc.redis != nil {
		data, _ := json.Marshal(users)
		if err := uc.redis.Set(ctx, cacheKey, data, 10*time.Minute).Err(); err != nil {
			span.LogFields(log.Error(err))
			fmt.Printf("⚠️ Gagal set Redis users: %v\n", err)
		}
	}

	return users, nil
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.GetUserByID")
	defer span.Finish()

	span.SetTag("operation", "GetUserByID")
	span.SetTag("user_id", id)

	result, err := uc.breaker.Execute(func() (interface{}, error) {
		return uc.userRepo.GetUserByID(ctx, id)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, err
	}
	return result.(*entity.User), nil
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user *entity.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.CreateUser")
	defer span.Finish()

	span.SetTag("operation", "CreateUser")
	span.SetTag("user_email", user.Email)

	if err := uc.validate.Struct(user); err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	_, err := uc.breaker.Execute(func() (interface{}, error) {
		return nil, uc.userRepo.CreateUser(ctx, user)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	if uc.redis != nil {
		if err := uc.redis.Del(ctx, "users:all").Err(); err != nil {
			span.LogFields(log.Error(err))
			fmt.Printf("⚠️ Gagal hapus cache users setelah Create: %v\n", err)
		}
	}

	return nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, id uint, user *entity.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.UpdateUser")
	defer span.Finish()

	span.SetTag("operation", "UpdateUser")
	span.SetTag("user_id", id)

	if err := uc.validate.Struct(user); err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	_, err := uc.breaker.Execute(func() (interface{}, error) {
		return nil, uc.userRepo.UpdateUser(ctx, id, user)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	if uc.redis != nil {
		if err := uc.redis.Del(ctx, "users:all").Err(); err != nil {
			span.LogFields(log.Error(err))
			fmt.Printf("⚠️ Gagal hapus cache users setelah Update: %v\n", err)
		}
	}

	return nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id uint) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.DeleteUser")
	defer span.Finish()

	span.SetTag("operation", "DeleteUser")
	span.SetTag("user_id", id)

	_, err := uc.breaker.Execute(func() (interface{}, error) {
		return nil, uc.userRepo.DeleteUser(ctx, id)
	})
	if err != nil {
		span.LogFields(log.Error(err))
		return err
	}

	if uc.redis != nil {
		if err := uc.redis.Del(ctx, "users:all").Err(); err != nil {
			span.LogFields(log.Error(err))
			fmt.Printf("⚠️ Gagal hapus cache users setelah Delete: %v\n", err)
		}
	}

	return nil
}
