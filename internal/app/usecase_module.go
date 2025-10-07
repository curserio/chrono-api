package app

import (
	"github.com/curserio/chrono-api/internal/usecase"
	"go.uber.org/fx"
)

var UsecaseModule = fx.Options(
	fx.Provide(
		usecase.NewMasterUseCase,
		usecase.NewServiceUseCase,
		usecase.NewBookingUseCase,
		usecase.NewClientUseCase,
		usecase.NewScheduleUseCase,
	),
)
