package converter

import (
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// ConvertPBGameToModelGame ...
func ConvertPBGameToModelGame(pbGame *registrator.Game) model.Game {
	return model.Game{
		ID:                  pbGame.GetId(),
		ExternalID:          pbGame.GetExternalId(),
		Type:                int32(pbGame.GetType()),
		Number:              pbGame.GetNumber(),
		Name:                pbGame.GetName(),
		Date:                model.DateTime(pbGame.GetDate().AsTime()),
		Price:               pbGame.GetPrice(),
		PaymentType:         pbGame.GetPaymentType(),
		MaxPlayers:          pbGame.GetMaxPlayers(),
		Payment:             model.PaymentType(pbGame.GetPayment()),
		Registered:          pbGame.GetRegistered(),
		My:                  pbGame.GetMy(),
		NumberOfMyLegioners: pbGame.GetNumberOfMyLegioners(),
		NumberOfLegioners:   pbGame.GetNumberOfLegioners(),
		NumberOfPlayers:     pbGame.GetNumberOfPlayers(),
		ResultPlace:         model.ResultPlace(pbGame.GetResultPlace()),
	}
}
