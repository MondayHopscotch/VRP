package routing

import (
	"vrp/internal/types"
)

// Solver is a general interface for route planners
type Solver interface {
	PlanRoutes() []types.Route
}
