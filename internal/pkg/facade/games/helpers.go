package games

import (
	"github.com/mono83/maybe"
	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

func convertProtoGameToModelGame(pbGame *gamepb.Game) model.Game {
	modelGame := model.Game{
		ID:          pbGame.GetId(),
		ExternalID:  maybe.Nothing[int32](),
		LeagueID:    int32(pbGame.GetLeagueId()),
		Type:        int32(pbGame.GetType()),
		Number:      pbGame.GetNumber(),
		Name:        maybe.Nothing[string](),
		PlaceID:     pbGame.GetPlaceId(),
		DateTime:    model.DateTime(pbGame.GetDate().AsTime()),
		Price:       pbGame.GetPrice(),
		PaymentType: maybe.Nothing[string](),
		MaxPlayers:  pbGame.GetMaxPlayers(),
		Payment:     maybe.Nothing[int32](),
		Registered:  pbGame.GetRegistered(),
		IsInMaster:  pbGame.GetIsInMaster(),
		HasPassed:   pbGame.GetHasPassed(),
		GameLink:    maybe.Nothing[string](),
	}

	if externalID := pbGame.GetExternalId(); externalID != nil {
		modelGame.ExternalID = maybe.Just(externalID.GetValue())
	}

	if name := pbGame.GetName(); name != nil {
		modelGame.Name = maybe.Just(name.GetValue())
	}

	if paymentType := pbGame.GetPaymentType(); paymentType != nil {
		modelGame.PaymentType = maybe.Just(paymentType.GetValue())
	}

	if payment := pbGame.Payment; payment != nil {
		modelGame.Payment = maybe.Just(int32(*payment))
	}

	if gameLink := pbGame.GetGameLink(); gameLink != nil {
		modelGame.GameLink = maybe.Just(gameLink.GetValue())
	}

	return modelGame
}
