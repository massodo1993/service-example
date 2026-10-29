package app

import (
	"context"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	authAPI "github.com/massodo1993/service-example/iam/internal/api/auth/v1"
	userAPI "github.com/massodo1993/service-example/iam/internal/api/user/v1"
	"github.com/massodo1993/service-example/iam/internal/config"
	"github.com/massodo1993/service-example/iam/internal/repository"
	sessionRepository "github.com/massodo1993/service-example/iam/internal/repository/session"
	userRepository "github.com/massodo1993/service-example/iam/internal/repository/user"
	"github.com/massodo1993/service-example/iam/internal/service"
	authService "github.com/massodo1993/service-example/iam/internal/service/auth"
	userService "github.com/massodo1993/service-example/iam/internal/service/user"
	"github.com/massodo1993/service-example/platform/closer"
	"github.com/massodo1993/service-example/platform/logger"
	"github.com/massodo1993/service-example/platform/pkg/cache"
	"github.com/massodo1993/service-example/platform/pkg/cache/redis"
	"github.com/massodo1993/service-example/platform/pkg/migrator/pg"
	authv1 "github.com/massodo1993/service-example/shared/pkg/proto/auth/v1"
	userv1 "github.com/massodo1993/service-example/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authV1API authv1.AuthServiceServer
	userV1API userv1.UserServiceServer

	authService service.AuthService
	userService service.UserService

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	pgPool *pgxpool.Pool

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthV1API(ctx context.Context) authv1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authAPI.NewAPI(d.AuthService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) UserV1API(ctx context.Context) userv1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userAPI.NewAPI(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(
			d.UserRepository(ctx),
			d.SessionRepository(),
			config.AppConfig().Session.TTL(),
		)
	}

	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PostgresPool(ctx))
	}

	return d.userRepository
}

func (d *diContainer) SessionRepository() repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisClient())
	}

	return d.sessionRepository
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		pool, err := pgxpool.New(connectCtx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to postgres: %s\n", err.Error()))
		}

		if err := pool.Ping(connectCtx); err != nil {
			panic(fmt.Sprintf("failed to ping postgres: %s\n", err.Error()))
		}

		migrationsDir := config.AppConfig().Postgres.MigrationsDir()
		migratorRunner := pg.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migrationsDir)

		if err := migratorRunner.Up(); err != nil {
			panic(fmt.Sprintf("failed to run migrations: %s\n", err.Error()))
		}

		closer.AddNamed("postgres pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}

		closer.AddNamed("redis pool", func(ctx context.Context) error {
			return d.redisPool.Close()
		})
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}
