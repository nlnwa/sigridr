package sigridr

import (
	log "github.com/sirupsen/logrus"
	"github.com/nu7hatch/gouuid"
)

func Uuid() string {
	id, err := uuid.NewV4()
	if err != nil {
		log.WithError(err).Fatal("Generating UUID")
	}
	return id.String()
}
