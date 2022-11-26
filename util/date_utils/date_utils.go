package date_utils

import (
	"time"

	"github.com/luizmoitinho/bookstore_users_api/config"
)

func GetNow() time.Time {
	loc, _ := time.LoadLocation(config.Propertie.LOCATION)
	return time.Now().In(loc)
}

func GetNowString() string {
	return GetNow().Format(config.Propertie.DATE_LAYOUT)
}
