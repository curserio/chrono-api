package app

import (
	"github.com/curserio/chrono-api/internal/handler"

	"go.uber.org/fx"
)

var HandlerModule = fx.Options(
	fx.Invoke(
		handler.NewMasterHandler,
		handler.NewServiceHandler,
		handler.NewBookingHandler,
		handler.NewClientHandler,
		handler.NewScheduleHandler,
	),
)
