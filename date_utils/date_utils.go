package date_utils

import "time"

const (
	API_DATE_LAYOUT = "2006-01-02T15:04:05Z"
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	return GetNow().Format(API_DATE_LAYOUT)
}
