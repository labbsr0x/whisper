package api

import "github.com/labbsr0x/whisper/web/config"
import "github.com/labbsr0x/whisper/db"

// AppAPI defines the available user apis
type AppAPI interface {
	Init(w *config.WebBuilder) AppAPI
}

// DefaultAppAPI holds the default implementation for the app api
type DefaultAppAPI struct {
	*config.WebBuilder
	dao db.AppDAO
}

var _dao = new(db.DefaultAppDAO)

// Init initializes the API
func (dapi *DefaultAppAPI) Init(w *config.WebBuilder) AppAPI {
	dapi.WebBuilder = w
	dapi.dao = _dao.Init(w.DB)

	return dapi
}
