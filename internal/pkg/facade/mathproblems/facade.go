//go:generate mockery --case underscore --name MathProblemServiceClient --with-expecter

package mathproblems

import (
	"context"

	mathproblempb "github.com/nikita5637/quiz-registrator-api/pkg/pb/math_problem"

	"google.golang.org/grpc"
)

// MathProblemServiceClient ...
type MathProblemServiceClient interface {
	SearchMathProblemByGameID(ctx context.Context, in *mathproblempb.SearchMathProblemByGameIDRequest, opts ...grpc.CallOption) (*mathproblempb.MathProblem, error)
}

// Facade ...
type Facade struct {
	mathProblemServiceClient MathProblemServiceClient
}

// Config ...
type Config struct {
	MathProblemServiceClient MathProblemServiceClient
}

// New ...
func New(cfg Config) *Facade {
	return &Facade{
		mathProblemServiceClient: cfg.MathProblemServiceClient,
	}
}
