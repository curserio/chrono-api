package app

import (
	"context"

	"github.com/curserio/chrono-api/config"
	"github.com/curserio/chrono-api/internal/repository"
	"github.com/curserio/chrono-api/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var RepositoryModule = fx.Options(
	fx.Provide(
		func(cfg *config.Config, lc fx.Lifecycle) (*pgxpool.Pool, error) {
			conn, err := postgres.InitConn(context.Background(), cfg.Database)
			if err != nil {
				return nil, err
			}

			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					conn.Close()
					return nil
				},
			})

			return conn, nil
		},
	),

	fx.Provide(
		postgres.NewMasterRepository,
		postgres.NewServiceRepository,
		postgres.NewBookingRepository,
		postgres.NewClientRepository,
		postgres.NewScheduleRepository,

		func(repo *postgres.MasterRepository) repository.MasterRepository {
			return repo
		},
		func(repo *postgres.ServiceRepository) repository.ServiceRepository {
			return repo
		},
		func(repo *postgres.BookingRepository) repository.BookingRepository {
			return repo
		},
		func(repo *postgres.ClientRepository) repository.ClientRepository {
			return repo
		},
		func(repo *postgres.ScheduleRepository) repository.ScheduleRepository {
			return repo
		},
	),
)
