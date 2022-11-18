package utils

import (
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

func Minute2Day(minute uint32) uint32 {
	t := time.Unix(int64(minute)*60, 0)
	return uint32(t.Year()*10000 + int(t.Month())*100 + t.Day())
}
func Second2Day(second uint32) uint32 {
	t := time.Unix(int64(second), 0)
	return uint32(t.Year()*10000 + int(t.Month())*100 + t.Day())
}
func AddDays(time uint32, days uint32) uint32 {
	sec := days * 24 * 60 * 60
	return time + sec
}
func Day2Second(date uint32) uint32 {
	if date == 0 {
		return 0
	}
	year := date / 10000
	month := (date % 10000) / 100
	day := date % 100
	if year == 0 || month == 0 || day == 0 {
		return 0
	}
	return uint32(time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Now().Location()).Unix())
}
func DiffTime(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()
	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()
	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)
	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}
func BeginTimeStampOfDate(t time.Time) int64 {
	timeStr := t.Format("2006-01-02")
	bt, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	timeStamp := bt.Unix()
	return timeStamp
}
func EndTimeStampOfDate(t time.Time) int64 {
	timeStr := t.Format("2006-01-02")
	bt, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	timeStamp := bt.Unix()
	return timeStamp + (24*60*60 - 1)
}
func BeginTimeStampOfYear(t time.Time) int64 {
	timeStr := t.Format("2006")
	bt, _ := time.ParseInLocation("2006", timeStr, time.Local)
	timeStamp := bt.Unix()
	return timeStamp
}
func EndTimeStampOfYear(t time.Time) int64 {
	timeStr := t.Format("2006")
	bt, _ := time.ParseInLocation("2006", timeStr, time.Local)
	timeStamp := bt.Unix()
	return timeStamp + (24*60*60*366 - 1)
}
func BeginTimeStampOfMonth(t time.Time) int64 {
	timeStr := t.Format("2006-01")
	bt, _ := time.ParseInLocation("2006-01", timeStr, time.Local)
	timeStamp := bt.Unix()
	return timeStamp
}
func EndTimeStampOfMonth(t time.Time) int64 {
	start := BeginTimeStampOfMonth(t)
	end := BeginTimeStampOfMonth(time.Unix(start+24*60*60*35, 0)) - 1
	return end
}
func UnixMillis() int64 {
	return time.Now().UnixNano() / 1e6
}
func UnixTimeStamp() int64 {
	return time.Now().Unix()
}
func UnixTimeStampUint32() uint32 {
	return uint32(fasthttp.CoarseTimeNow().Unix())
}
func Now() uint32 {
	return uint32(fasthttp.CoarseTimeNow().Unix())
}
func NowMs() uint64 {
	now := time.Now()
	return uint64(now.Unix()*1000) + (uint64(now.Nanosecond()) / uint64(time.Millisecond))
}
func DateYmd() string {
	return time.Now().Format("2006-01-02")
}
func TimeNowCompactYmd() string {
	return time.Now().Format("20060102")
}
func DateMMDD() string {
	return time.Now().Format("0102")
}
func TimeNowhhmmss() string {
	return time.Now().Format("150405")
}
func DateYmdhis() string {
	return time.Now().Format("2006年1月2日 15:04:05")
}
func DateYmdhi() string {
	return time.Now().Format("2006年1月2日 15:04")
}
func TimeStamp2DateYmd(ts uint32) string {
	return time.Unix(int64(ts), 0).Format("2006-01-02")
}
func TimeStamp2DateYmdhis(ts uint32) string {
	return time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
}
func TimeStamp2DateYmdhi(ts uint32) string {
	return time.Unix(int64(ts), 0).Format("2006-01-02 15:04")
}
func TimeStamp2DateCompactYmdhis(ts uint32) string {
	return time.Unix(int64(ts), 0).Format("20060102150405")
}
func RFC3339Nano2DateYmdhis(t string) string {
	// RFC3339Nano type define at time/format.go
	// e.g: 2020-08-21T16:21:11+08:00
	ts, _ := time.Parse(time.RFC3339Nano, t)
	return ts.Format("2006-01-02 15:04:05")
}
func RFC3339Nano2DateYmdhi(t string) string {
	ts, _ := time.Parse(time.RFC3339Nano, t)
	return ts.Format("2006-01-02 15:04")
}
func RFC3339Nano2TimeStamp(t string) int64 {
	ts, _ := time.Parse(time.RFC3339Nano, t)
	return ts.Unix()
}
func RFC33392TimeStamp(t string) int64 {
	ts, _ := time.Parse(time.RFC3339, t)
	return ts.Unix()
}
func BeginTimeOfDate(t time.Time) time.Time {
	timeStr := t.Format("2006-01-02")
	bt, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return bt
}
func EndTimeOfDate(t time.Time) time.Time {
	timeStr := t.Format("2006-01-02")
	bt, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return bt.Add(24*time.Hour - 1)
}
func StringToTimestamp(date string) uint32 {
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", date, loc)
	return uint32(t.Unix())
}

// 获取需要统计的日期 以及是否包含当天
func GetDateListBetween2Times(startTime, endTime int64) (out []uint32, haveToday bool) {
	if startTime > endTime {
		return nil, false
	}
	st := time.Unix(startTime, 0)
	endDate64, _ := strconv.ParseUint(time.Unix(endTime, 0).Format("20060102"), 10, 64)
	endDate := uint32(endDate64)

	startDate := uint32(0)
	todayDate := GetTodayDate()
	for startDate != endDate {

		startDate64, _ := strconv.ParseUint(st.Format("20060102"), 10, 64)
		startDate := uint32(startDate64)
		out = append(out, startDate)
		st = st.Add(time.Second * 3600 * 24)
		if startDate >= todayDate {
			haveToday = true
			return
		}

		if startDate >= endDate {
			return
		}
	}
	return
}

func GetTodayDate() uint32 {
	today64, _ := strconv.ParseUint(time.Now().Format("20060102"), 10, 64)
	return uint32(today64)
}
