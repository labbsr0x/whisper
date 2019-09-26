package workers

import (
	"github.com/labbsr0x/whisper/web/config"
)

var Mail DefaultMailWorkerAPI

// InitFromWebBuilder initializes the workers
func InitFromWebBuilder (builder *config.WebBuilder) {
	Mail.InitFromWebBuilder(builder)
}

// Run start all the workers
func Run () {
	Mail.Start()
}