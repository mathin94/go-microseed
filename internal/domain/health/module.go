package health

import (
	"microseed/internal/httpx"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewHandler,
			fx.As(new(httpx.RouteRegistrar)),
			fx.ResultTags(`group:"routes"`),
		),
	),
)
