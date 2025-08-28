package api

import (
	"net/http"
	"time"
)

func nextDateHandler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeat := query.Get("repeat")

	if dateStr == "" || repeat == "" {
		http.Error(w, "require query parameters 'date' and 'repeat'", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		t, err := time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "incorrect 'now' parameter", http.StatusBadRequest)
			return
		}
		now = t
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(next))
}
