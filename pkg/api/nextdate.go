package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("incorrect format for repeat rule")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("incorrect amount of days in repeat")
		}

		if days < 1 || days > 400 {
			return "", fmt.Errorf("the number of days must be between 1 and 400")
		}
		date = date.AddDate(0, 0, days)
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, days)
		}
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("incorrect format for repeat rule")
		}

		date = date.AddDate(1, 0, 0)

		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}
	default:
		return "", fmt.Errorf("incorrect repeat format:%s", parts[0])
	}
	return date.Format(DateFormat), nil
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "incorrect now value", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

func afterNow(date time.Time, now time.Time) bool {
	return date.After(now)
}
