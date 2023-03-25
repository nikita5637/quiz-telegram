package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DictLength(t *testing.T) {
	assert.Len(t, dictionary, 1)
	assert.Len(t, dictionary["ru"], 70)
}

func Test_DictPhrases(t *testing.T) {
	assert.Equal(t, `Список доступных команд:
/games - список всех игр
/mygames - список игр, на которые ты идёшь
/photos - фотографии с игр
/registeredgames - список игр, на которые мы зарегистрированы
/settings - настройки
/help - помощь

ℹ️ - ты идёшь на эту игру
❗️ - ты не идёшь на эту игру, но идут другие

Формат игроков "игроки/легионеры/всего может быть игроков"`, dictionary["ru"]["help_message"])

	assert.Equal(t, `Привет, %s!
Я Зоя :)
Добро пожаловать :)
Список доступных команд:
/games - список всех игр
/mygames - список игр, на которые ты идёшь
/photos - фотографии с игр
/registeredgames - список игр, на которые мы зарегистрированы
/settings - настройки
/help - помощь`, dictionary["ru"]["welcome_message"])
}
