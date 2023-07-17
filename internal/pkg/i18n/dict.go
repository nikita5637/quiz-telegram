package i18n

import "context"

const defaultLang = "ru"

func translate(ctx context.Context, key string, defaultString string) string {
	lang := GetLangFromContext(ctx)

	v, ok := dictionary[lang][key]
	if ok {
		return v
	}
	return defaultString
}

var dictionary = map[string]map[string]string{
	"ru": {
		"address":                              "Адрес",
		"add_to_calendar":                      "Добавить в календарь",
		"birthdate_changed":                    "Дата рождения изменена",
		"card":                                 "Карта",
		"cash":                                 "Наличные",
		"cash_game_payment":                    "Играем за деньги",
		"certificate":                          "Сертификат",
		"change_birthdate":                     "Изменить дату рождения",
		"change_email":                         "Изменить email",
		"change_name":                          "Изменить имя",
		"change_phone":                         "Изменить номер телефона",
		"change_sex":                           "Изменить пол",
		"choose_a_game":                        "Выбери игру",
		"datetime":                             "Дата и время",
		"email_changed":                        "Email изменён",
		"enter_your_birthdate":                 "Окей. Введи свою дату рождения в формате ДД.ММ.ГГГГ.",
		"enter_your_email":                     "Окей. Введи свой email.",
		"enter_your_name":                      "Окей. Введи своё имя.",
		"enter_your_phone":                     "Окей. Введи свой номер телефона в формате +79XXXXXXXXXX.",
		"enter_your_sex":                       "Окей. Введи свой пол(М или Ж)",
		"free_game_payment":                    "Играем по сертификату",
		"game_cost":                            "Стоимость игры",
		"game_not_found":                       "Игра не найдена",
		"game_photos":                          "Фотографии с игр",
		"help_message":                         "Список доступных команд:\n/games - список всех игр\n/mygames - список игр, на которые ты идёшь\n/photos - фотографии с игр\n/registeredgames - список игр, на которые мы зарегистрированы\n/settings - настройки\n/help - помощь\n\nℹ️ - ты идёшь на эту игру\n❗️ - ты не идёшь на эту игру, но идут другие\n\nФормат игроков \"игроки/легионеры/всего может быть игроков\"",
		"legioner_by":                          "Лег от",
		"legioner_is_likely_to_come":           "Лег придёт",
		"legioner_is_signed_up_for_the_game":   "Легионер записан на игру",
		"legioner_is_unsigned_up_for_the_game": "Легионер выписан из игры",
		"legioner_is_unlikely_to_come":         "Лег мб придёт",
		"legioner_will_not_come":               "Лег не придёт",
		"list_of_games_is_empty":               "Нет игр",
		"list_of_games_with_photos_is_empty":   "Нет игр с фотографиями",
		"list_of_my_games_is_empty":            "Ты пока не играешь с нами",
		"list_of_players":                      "Список игроков",
		"list_of_players_is_empty":             "Нет игроков",
		"list_of_registered_games":             "Список зарегистрированных игр",
		"list_of_registered_games_is_empty":    "Нет зарегистрированных игр",
		"list_of_your_games":                   "Список твоих игр",
		"mix":                                  "Микс",
		"mix_game_payment":                     "Кто-то платит, кто-то нет",
		"my_games":                             "Мои игры",
		"name_changed":                         "Имя изменено",
		"no_free_slot":                         "Нет свободных мест",
		"number":                               "Номер",
		"number_of_players":                    "Количество игроков",
		"payment":                              "Оплата",
		"permission_denied":                    "Доступ запрещён",
		"phone_changed":                        "Номер телефона изменён",
		"place":                                "Место",
		"player_is_likely_to_come":             "Я точно приду",
		"player_is_unlikely_to_come":           "Я мб приду",
		"player_will_not_come":                 "Я не приду",
		"plays_likely":                         "Играет",
		"plays_unlikely":                       "Под вопросом",
		"register_for_lottery":                 "Зарегистрироваться в лотерее",
		"registered_game":                      "Мы зарегистрированы на игру",
		"registered_games":                     "Зарегистрированные игры",
		"registration_for_a_game":              "Регистрация на игры",
		"registration_link":                    "Ссылка на регистрацию",
		"remind_that_there_is_a_game_today":    "Напоминаю, что сегодня играем",
		"remind_that_there_is_a_lottery":       "Напоминаю о регистрации в лотерее",
		"route_to_bar":                         "Маршрут до бара",
		"settings":                             "Настройки",
		"sex_changed":                          "Пол изменён",
		"something_went_wrong":                 "Что-то пошло не так",
		"time":                                 "Время",
		"title":                                "Название",
		"unregistered_game":                    "Мы не зарегистрированы на игру",
		"welcome_message":                      "Привет, %s!\nЯ Зоя :)\nДобро пожаловать :)\nСписок доступных команд:\n/games - список всех игр\n/mygames - список игр, на которые ты идёшь\n/photos - фотографии с игр\n/registeredgames - список игр, на которые мы зарегистрированы\n/settings - настройки\n/help - помощь",
		"you_are_already_registered_for_the_game":         "Ты уже зарегистрирован на игру",
		"you_are_registered_for_the_game":                 "Ты зарегистрирован на игру",
		"you_are_signed_up_for_the_game":                  "Ты записан на игру",
		"you_are_unsigned_up_for_the_game":                "Ты выписан из игры",
		"you_have_successfully_registered_in_the_lottery": "Ты зарегистрирован в лотерее",
		"your_lottery_number_is":                          "Твой лотерейный номер",
		"zoya":                                            "Зоя",
	},
}
