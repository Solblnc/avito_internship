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
		jsonRespond(w, http.StatusInternalServerError, []byte(`server error`))
		return
	}

	result := make([]string, 0)

	result = s.DB.GetActualSegments(userId)

	res, err := json.Marshal(result)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`server error`))
		return
	}
	jsonRespond(w, 200, res)
}
