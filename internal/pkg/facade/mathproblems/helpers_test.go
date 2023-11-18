package mathproblems

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/mathproblems/mocks"
)

type fixture struct {
	ctx context.Context

	mathProblemServiceClient *mocks.MathProblemServiceClient

	facade *Facade
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		mathProblemServiceClient: mocks.NewMathProblemServiceClient(t),
	}

	fx.facade = New(Config{
		MathProblemServiceClient: fx.mathProblemServiceClient,
	})

	t.Cleanup(func() {})

	return fx
}
