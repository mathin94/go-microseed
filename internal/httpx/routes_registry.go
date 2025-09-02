package httpx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouteRegistrar interface {
	Register(r *gin.Engine)
}

type routesIn struct {
	fx.In

	Engine     *gin.Engine
	Registrars []RouteRegistrar `group:"routes"`
}

func registerAll(in routesIn) {
	for _, rr := range in.Registrars {
		rr.Register(in.Engine)
	}
}

var RoutesModule = fx.Options(
	fx.Invoke(registerAll),
)
