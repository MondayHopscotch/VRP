package routing

import (
	"vrp/internal/types"
)

type Solver interface {
	PlanRoutes() []types.Route
}
