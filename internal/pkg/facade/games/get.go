package games

import (
	"context"
	"fmt"
	"sort"

	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetGame ...
func (f *Facade) GetGame(ctx context.Context, id int32) (model.Game, error) {
	pbGame, err := f.gameServiceClient.GetGame(ctx, &gamepb.GetGameRequest{
		Id: id,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return model.Game{}, ErrGameNotFound
		}

		return model.Game{}, fmt.Errorf("getting game error: %w", err)
	}

	return convertProtoGameToModelGame(pbGame), nil
}

// GetGames ...
func (f *Facade) GetGames(ctx context.Context, registered, isInMaster, hasPassed bool) ([]model.Game, error) {
	gamesResp, err := f.gameServiceClient.ListGames(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("get games error: %w", err)
	}

	modelGames := make([]model.Game, 0, len(gamesResp.GetGames()))
	for _, pbGame := range gamesResp.GetGames() {
		if pbGame.GetRegistered() == registered && pbGame.GetIsInMaster() == isInMaster && pbGame.GetHasPassed() == hasPassed {
			for _, t := range f.permittedGameTypes {
				if gamepb.GameType(t) == pbGame.GetType() {
					modelGames = append(modelGames, convertProtoGameToModelGame(pbGame))
					break
				}
			}
		}
	}

	return modelGames, nil
}

// GetGamesByUserID ...
func (f *Facade) GetGamesByUserID(ctx context.Context, userID int32) ([]model.Game, error) {
	userGameIDs, err := f.gamePlayersFacade.GetUserGameIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting user game IDs error: %w", err)
	}

	batchGetGamesResp, err := f.gameServiceClient.BatchGetGames(ctx, &gamepb.BatchGetGamesRequest{
		Ids: userGameIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("batch getting games error: %w", err)
	}

	modelGames := make([]model.Game, 0, len(batchGetGamesResp.GetGames()))
	for _, pbGame := range batchGetGamesResp.GetGames() {
		if !pbGame.GetHasPassed() {
			modelGames = append(modelGames, convertProtoGameToModelGame(pbGame))
		}
	}

	sort.Slice(modelGames, func(i, j int) bool {
		return modelGames[i].DateTime.AsTime().Before(modelGames[j].DateTime.AsTime())
	})

	return modelGames, nil
}
