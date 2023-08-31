package server

import (
	_ "avito_internship/docs"
	"encoding/json"
	"net/http"
	"strconv"
)

// @Summary	GetSegments
// @Description	get segments for specific user by userId
// @Tags user
// @Param userId path int true "User id"
// @Success 200	"Segments for user:"
// @Router	/user/get_segments{id} [get]
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

	result = s.Service.GetActualSegments(userId)

	res, err := json.Marshal(result)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`app error`))
		return
	}
	jsonRespond(w, 200, res)
}

// @Summary		get user segments history
// @Description	get user segments history for year and month
// @Tags			user
// @Param			year	path		int					true	"year"
// @Param			month	path		int					true	"month"
// @Success		200		"csv file completed"
// @Router			/segment/get_history{year}{month} [get]
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

	result, err := s.Service.GetHistory(year, month)
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

// @Summary		create users
// @Description	create test users just for testing an app
// @Tags			user
// @Success		200		"users are created"
// @Router			/user/create_user [get]
func (s *Server) CreateUsers(w http.ResponseWriter, r *http.Request) {
	err := s.Service.CreateUser()
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot create users`))
		return
	}
	jsonRespond(w, http.StatusOK, []byte(`users are created`))
}
