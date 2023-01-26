package games

import (
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(err error) error {
	if err == nil {
		return nil
	}

	st := status.Convert(err)
	if st.Code() == codes.NotFound {
		return model.ErrGameNotFound
	}

	return err

}
