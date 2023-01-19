package request

import (
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/nikita5637/quiz-telegram/internal/pkg/storage"
)

// Facade ...
type Facade struct {
	cache          map[string]model.Request
	requestStorage storage.RequestStorage
}

// Config ...
type Config struct {
	RequestStorage storage.RequestStorage
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		cache:          make(map[string]model.Request),
		requestStorage: cfg.RequestStorage,
	}
}
