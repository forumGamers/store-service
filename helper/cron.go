package helper

import "time"

func SetSchedule24() time.Time {
	now := time.Now()

	return now.Add(24 * time.Hour)
}

//contoh
// _, err := c.AddFunc(fmt.Sprintf("%d %d %d %d %d *", after24Hours.Second(), after24Hours.Minute(), after24Hours.Hour(), after24Hours.Day(), after24Hours.Month()), func() {
// 	checkTransaction(transaction.ID)
// })