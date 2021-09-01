package stop

import (
	"fmt"
	"time"
)

//TimestampToDate Timestamp(miliseconds) -> 02/01/2006 15:04:05
func TimestampToDate(timestamp int64) string {
	t := time.Unix(0, timestamp*int64(time.Millisecond))
	strDate := t.Format("02/01/2006 15:04:05")
	return strDate
}

//TimestampToDate Timestamp(miliseconds) -> 02-01-2006
func TimestampToDate1(timestamp int64) string {
	t := time.Unix(0, timestamp*int64(time.Millisecond))
	strDate := t.Format("02-01-2006")
	return strDate
}

//TimestampToDateFile Timestamp(miliseconds) -> 2006-01-02 15:04:05
func TimestampToDateFile(timestamp int64) string {
	t := time.Unix(0, timestamp*int64(time.Millisecond))
	fmt.Println(t)
	strDate := t.Format("2006-01-02 15:04:05")
	return strDate
}

//DateToTimestamp 2006-01-02 15:04:05 -> Timestamp(miliseconds)
func DateToTimestamp(date string) int64 {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, date)
	t1 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	if err != nil {
		fmt.Println(err)
	}
	back := t1.UnixNano() / int64(time.Millisecond)
	// fmt.Print(t)
	// fmt.Print(" ")
	// fmt.Println(back)
	return back
}

//DateToTimestampRange 2006-01-02 -> Timestamp(miliseconds)
func DateToTimestampRange(date string) int64 {
	layout := "2006-01-02"
	t, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println(err)
	}
	return t.UnixNano() / int64(time.Millisecond)
}
