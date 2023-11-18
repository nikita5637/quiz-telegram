package mathproblems

import (
	"errors"
	"testing"

	mathproblempb "github.com/nikita5637/quiz-registrator-api/pkg/pb/math_problem"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFacade_GetMathProblemByGameID(t *testing.T) {
	t.Run("error: math problem not found", func(t *testing.T) {
		fx := tearUp(t)

		fx.mathProblemServiceClient.EXPECT().SearchMathProblemByGameID(fx.ctx, &mathproblempb.SearchMathProblemByGameIDRequest{
			GameId: 1,
		}).Return(nil, status.New(codes.NotFound, "some error").Err())

		got, err := fx.facade.GetMathProblemByGameID(fx.ctx, 1)
		assert.Equal(t, model.MathProblem{}, got)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrMathProblemNotFound)
	})

	t.Run("error: some error", func(t *testing.T) {
		fx := tearUp(t)

		fx.mathProblemServiceClient.EXPECT().SearchMathProblemByGameID(fx.ctx, &mathproblempb.SearchMathProblemByGameIDRequest{
			GameId: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetMathProblemByGameID(fx.ctx, 1)
		assert.Equal(t, model.MathProblem{}, got)
		assert.Error(t, err)

		st := status.Convert(err)
		assert.Equal(t, codes.Unknown, st.Code())
		assert.Len(t, st.Details(), 0)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.mathProblemServiceClient.EXPECT().SearchMathProblemByGameID(fx.ctx, &mathproblempb.SearchMathProblemByGameIDRequest{
			GameId: 1,
		}).Return(&mathproblempb.MathProblem{
			Id:     1,
			GameId: 1,
			Url:    "url",
		}, nil)

		got, err := fx.facade.GetMathProblemByGameID(fx.ctx, 1)
		assert.Equal(t, model.MathProblem{
			ID:     1,
			GameID: 1,
			URL:    "url",
		}, got)
		assert.NoError(t, err)
	})
}
