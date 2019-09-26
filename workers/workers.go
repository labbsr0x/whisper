package workers

import (
	"github.com/labbsr0x/whisper/web/config"
)

type Workers struct {
	*config.WebBuilder
	Mail     MailWorkerAPI
}

func (w *Workers) InitFromWebBuilder (builder *config.WebBuilder) *Workers {
	w.Mail = new(DefaultMailWorkerAPI).InitFromWebBuilder(builder)

	return w
}