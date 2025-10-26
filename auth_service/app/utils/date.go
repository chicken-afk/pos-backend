package utils

import "time"

func ConvertEpochToDateTimeJakarta(epoch int64) string {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return ""
	}
	return ConvertEpochToDateTimeInLocation(epoch, location)
}

func ConvertEpochToDateTimeInLocation(epoch int64, location *time.Location) string {
	t := time.Unix(epoch, 0).In(location)
	return t.Format("2006-01-02T15:04:05+07:00")
}

func ConvertEpochToDuration(epoch int64) time.Duration {
	currentTime := time.Now().Unix()
	durationInSeconds := epoch - currentTime
	return time.Duration(durationInSeconds) * time.Second
}
