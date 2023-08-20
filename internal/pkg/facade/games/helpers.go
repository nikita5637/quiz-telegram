package games

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(err error) error {
	if err == nil {
		return nil
	}

	st := status.Convert(err)
	if st.Code() == codes.NotFound {
		return ErrGameNotFound
	}

	return err

}
