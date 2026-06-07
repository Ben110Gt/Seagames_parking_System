package util

import (
	"math"
	"os"
	"strconv"
	"time"
)

func CalcFine(checkin, checkout time.Time) (int64, int) {
	rateKip, _ := strconv.ParseInt(os.Getenv("DAILY_FINE_KIP"), 10, 64)
	if rateKip == 0 {
		rateKip = 2000
	}

	loc, _ := time.LoadLocation("Asia/Vientiane")
	inDay := checkin.In(loc).Truncate(24 * time.Hour)
	outDay := checkout.In(loc).Truncate(24 * time.Hour)

	days := int(math.Round(outDay.Sub(inDay).Hours() / 24))
	if days <= 1 {
		return 0, days
	}
	extraDays := days - 1
	return int64(extraDays) * rateKip, days
}
