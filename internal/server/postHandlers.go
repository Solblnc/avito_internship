package server

import (
	"avito_internship/internal/model"
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) CreateSegment(w http.ResponseWriter, r *http.Request) {
	var segment model.Segment
	data, err := io.ReadAll(r.Body)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`error in reading a body`))
		return
	}

	if err = json.Unmarshal(data, &segment); err != nil {
		jsonRespond(w, http.StatusBadRequest, []byte(`error in unmarshalling a body`))
		return
	}
	_, err = s.DB.Create(segment.Name)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot create a segment`))
		return
	}

	jsonRespond(w, http.StatusOK, []byte(`segment is created`))

}

func (s *Server) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	var segment model.Segment
	data, err := io.ReadAll(r.Body)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`error in reading a body`))
		return
	}

	if err = json.Unmarshal(data, &segment); err != nil {
		jsonRespond(w, http.StatusBadRequest, []byte(`error in unmarshalling a body`))
		return
	}
	err = s.DB.Delete(segment.Name)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot delete a segment`))
		return
	}

	jsonRespond(w, http.StatusOK, []byte(`segment is deleted`))

}

func (s *Server) AddUserToSegment(w http.ResponseWriter, r *http.Request) {
	var input model.InputAddUser
	data, err := io.ReadAll(r.Body)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`error in reading a body`))
		return
	}

	if err = json.Unmarshal(data, &input); err != nil {
		jsonRespond(w, http.StatusBadRequest, []byte(`error in unmarshalling a body`))
		return
	}
	err = s.DB.AddUser(input.SegmentsAdd, input.SegmentsDelete, input.UserId)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot add/delete segments to user`))
		return
	}

	jsonRespond(w, http.StatusOK, []byte(`segments are added/deleted to user`))
}
