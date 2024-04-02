package routing

import (
	"vrp/internal"
)

type Solver interface {
	PlanRoutes() []internal.Route
}
