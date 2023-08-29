package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Server) GetSegments(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()
	Id, ok := m["user_id"]
	if !ok {
		jsonRespond(w, http.StatusBadRequest, []byte(`invalid request`))
		return
	}

	userId, err := strconv.Atoi(Id[0])
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`app error`))
		return
	}

	result := make([]string, 0)

	result = s.DB.GetActualSegments(userId)

	res, err := json.Marshal(result)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`app error`))
		return
	}
	jsonRespond(w, 200, res)
}

func (s *Server) GetHistory(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()
	Year, ok := m["year"]
	Month, ok := m["month"]
	if !ok {
		jsonRespond(w, http.StatusBadRequest, []byte(`invalid request`))
		return
	}

	year, err := strconv.Atoi(Year[0])
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`app error`))
		return
	}

	month, err := strconv.Atoi(Month[0])
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`app error`))
		return
	}

	result, err := s.DB.GetHistory(year, month)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot get history from database`))
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=history.csv")

	_, err = w.Write([]byte(result))
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot send a csv file`))
		return
	}

	jsonRespond(w, 200, []byte(`csv file completed`))
}
