package date_utils

import (
	"time"

	"github.com/luizmoitinho/bookstore_users_api/config"
)

func GetNow() time.Time {
	loc, _ := time.LoadLocation(config.Propertie.LOCATION)
	return time.Now().In(loc)
}

func GetNowApiFormat() string {
	return GetNow().Format(config.Propertie.API_DATE_LAYOUT)
}

func GetNowDbFormat() string {
	return GetNow().Format(config.Propertie.DB_DATE_LAYOUT)
}
