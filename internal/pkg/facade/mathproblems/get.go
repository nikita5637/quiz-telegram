package mathproblems

import (
	"context"

	mathproblem "github.com/nikita5637/quiz-registrator-api/pkg/pb/math_problem"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetMathProblemByGameID ...
func (f *Facade) GetMathProblemByGameID(ctx context.Context, gameID int32) (model.MathProblem, error) {
	pbMathProblem, err := f.mathProblemServiceClient.SearchMathProblemByGameID(ctx, &mathproblem.SearchMathProblemByGameIDRequest{
		GameId: gameID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return model.MathProblem{}, ErrMathProblemNotFound
		}

		return model.MathProblem{}, err
	}

	return model.MathProblem{
		ID:     pbMathProblem.GetId(),
		GameID: pbMathProblem.GetGameId(),
		URL:    pbMathProblem.GetUrl(),
	}, nil
}
