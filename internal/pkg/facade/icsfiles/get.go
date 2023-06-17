package icsfiles

import (
	"context"

	icsfilemanagerpb "github.com/nikita5637/quiz-ics-manager-api/pkg/pb/ics_file_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// GetICSFileByGameID ...
func (f *Facade) GetICSFileByGameID(ctx context.Context, gameID int32) (model.ICSFile, error) {
	resp, err := f.icsFileManagerAPIClient.GetICSFileByGameID(ctx, &icsfilemanagerpb.GetICSFileByGameIDRequest{
		GameId: gameID,
	})
	if err != nil {
		return model.ICSFile{}, err
	}

	return model.ICSFile{
		ID:     resp.GetId(),
		GameID: resp.GetGameId(),
		Name:   resp.GetName(),
	}, nil
}
