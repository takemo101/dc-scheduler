package database

import (
	"github.com/takemo101/dc-scheduler/pkg/infrastructure"
)

// Models gorm model list
var Models = []interface{}{
	&infrastructure.Admin{},
	&infrastructure.Bot{},
	&infrastructure.PostMessage{},
	&infrastructure.SentMessage{},
	&infrastructure.ScheduleTiming{},
}