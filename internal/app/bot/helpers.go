package bot

import (
	"time"

	"github.com/mono83/maybe"
	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

var (
	listOfPlayersLexeme = i18n.Lexeme{
		Key:      "list_of_players",
		FallBack: "List of players",
	}
	unknownLexeme = i18n.Lexeme{
		Key:      "unknown",
		FallBack: "Unknown",
	}
	youAreAlreadyRegisteredForTheGameLexeme = i18n.Lexeme{
		Key:      "you_are_already_registered_for_the_game",
		FallBack: "You are already registered for the game",
	}
	youAreSignedUpForTheGameLexeme = i18n.Lexeme{
		Key:      "you_are_signed_up_for_the_game",
		FallBack: "You are signed up for the game",
	}
)

func convertProtoUserToModelUser(pbUser *usermanagerpb.User) model.User {
	modelUser := model.User{
		ID:         pbUser.GetId(),
		Name:       pbUser.GetName(),
		TelegramID: pbUser.GetTelegramId(),
		Email:      maybe.Nothing[string](),
		Phone:      maybe.Nothing[string](),
		State:      int32(pbUser.GetState()),
		Birthdate:  maybe.Nothing[string](),
		Sex:        maybe.Nothing[model.Sex](),
	}

	if email := pbUser.GetEmail(); email != nil {
		modelUser.Email = maybe.Just(email.GetValue())
	}

	if phone := pbUser.GetPhone(); phone != nil {
		modelUser.Phone = maybe.Just(phone.GetValue())
	}

	if birthdate := pbUser.GetBirthdate(); birthdate != nil {
		birthdateTime, err := time.Parse("2006-01-02", birthdate.GetValue())
		if err == nil {
			modelUser.Birthdate = maybe.Just(birthdateTime.Format("02.01.2006"))
		}
	}

	if pbUser != nil && pbUser.Sex != nil {
		modelUser.Sex = maybe.Just(model.Sex(pbUser.GetSex()))
	}

	return modelUser
}
