package i18n

import "context"

const defaultLang = "ru"

// Translate ...
func Translate(ctx context.Context, key string, defaultString string) string {
	lang := GetLangFromContext(ctx)

	v, ok := dictionary[lang][key]
	if ok {
		return v
	}
	return defaultString
}

var dictionary = map[string]map[string]string{
	"ru": {
		"back_to_the_games_list":                  "Назад к списку игр",
		"choose_a_game":                           "Выбери игру",
		"cash_game_payment":                       "Играем за денюжку",
		"change_email":                            "Изменить email",
		"change_name":                             "Изменить имя",
		"change_phone":                            "Изменить номер телефона",
		"email_changed":                           "Email изменён",
		"enter_your_email":                        "Окей. Введи свой email.",
		"enter_your_name":                         "Окей. Введи своё имя.",
		"enter_your_phone":                        "Окей. Введи свой номер телефона в формате +79XXXXXXXXXX.",
		"free_game_payment":                       "Играем по сертификату",
		"game_not_found":                          "Игра не найдена",
		"game_photos":                             "Фотографии с игр",
		"help_message":                            "Список доступных команд:\n/games - список всех игр\n/mygames - список игр, на которые ты идёшь\n/photos - фотографии с игр\n/registeredgames - список игр, на которые мы зарегистрированы\n/settings - настройки\n/help - помощь\n\nℹ️ - ты идёшь на эту игру\n❗️ - ты не идёшь на эту игру, но идут другие\n\nФормат игроков \"игроки/легионеры/всего может быть игроков\"",
		"legioner_by":                             "Лег от",
		"legioner_is_likely_to_come":              "Лег придёт",
		"legioner_is_unlikely_to_come":            "Лег может быть придёт",
		"legioner_will_not_come":                  "Лег не придёт",
		"list_of_games_is_empty":                  "Нет игр",
		"list_of_games_with_photos_is_empty":      "Нет игр с фотографиями",
		"list_of_my_games_is_empty":               "Ты пока не играешь с нами",
		"list_of_players":                         "Список игроков",
		"list_of_players_is_empty":                "Нет игроков",
		"list_of_registered_games":                "Список зарегистрированных игр",
		"list_of_registered_games_is_empty":       "Нет зарегистрированных игр",
		"list_of_your_games":                      "Список твоих игр",
		"mix_game_payment":                        "Кто-то платит, кто-то нет",
		"name_changed":                            "Имя изменено",
		"no_free_slot":                            "Нет свободных мест",
		"permission_denied":                       "Доступ запрещён",
		"phone_changed":                           "Номер телефона изменён",
		"place":                                   "Место",
		"player_is_likely_to_come":                "Я точно приду",
		"player_is_unlikely_to_come":              "Я может быть приду",
		"player_will_not_come":                    "Я не приду",
		"plays_likely":                            "Играет",
		"plays_unlikely":                          "Под вопросом",
		"register_for_lottery":                    "Зарегистрироваться в лотерее",
		"registered_game":                         "Мы зарегистрированы на игру",
		"registration_for_a_game":                 "Регистрация на игры",
		"remind_that_there_is_a_game_today":       "Напоминаю, что сегодня играем",
		"route_to_bar":                            "Маршрут до бара",
		"settings":                                "Настройки",
		"something_went_wrong":                    "Что-то пошло не так",
		"time":                                    "Время",
		"unregistered_game":                       "Мы не зарегистрированы на игру",
		"welcome_message":                         "Привет, %s!\nЯ Зоя :)\nДобро пожаловать :)\nСписок доступных команд:\n/games - список всех игр\n/mygames - список игр, на которые ты идёшь\n/photos - фотографии с игр\n/registeredgames - список игр, на которые мы зарегистрированы\n/settings - настройки\n/help - помощь",
		"you_are_already_registered_for_the_game": "Ты уже зарегистрирован на игру",
		"you_are_registered_for_the_game":         "Ты зарегистрирован на игру",
		"you_have_successfully_registered_in_the_lottery": "Ты зарегистрирован в лотерее",
		"your_lottery_number_is":                          "Твой лотерейный номер",
		"zoya":                                            "Зоя",
	},
}
