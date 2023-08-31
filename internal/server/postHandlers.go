package server

import (
	_ "avito_internship/docs"
	"avito_internship/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// @Summary		create segments
// @Description	create segment and randomly add users
// @Tags			segment
// @Param			segmentName	path		string					true	"segment name"
// @Param			percent 	path		int 					true	"percentage"
// @Success		200		"segment is created"
// @Router			/user/create_segment [post]
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
	_, err = s.Service.Create(segment.Name, segment.Percent)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot create a segment`))
		return
	}

	jsonRespond(w, http.StatusOK, []byte(`segment is created`))

}

// @Summary		delete segment
// @Description	delete segment
// @Tags			segment
// @Param			segmentName	path		string					true	"segment name"
// @Success		200		"segment is deleted"
// @Router			/user/delete_segment [delete]
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
	err = s.Service.Delete(segment.Name)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot delete a segment`))
		return
	}

	jsonRespond(w, http.StatusOK, []byte(`segment is deleted`))

}

// @Summary		add/delete segments to user
// @Description	adding and deleting specific segments for specific user
// @Tags			user
// @Param			segmentAdd 	path		string					true	"segment name"
// @Param			segmentDelete 	path		string					true	"segment name"
// @Param			userId	path		int					true	"User id"
// @Success		200		"Segments for user:"
// @Router			/segment/add_user_segment [post]
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
	err = s.Service.AddUser(input.SegmentsAdd, input.SegmentsDelete, input.UserId)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot add/delete segments to user`))
		return
	}

	jsonRespond(w, http.StatusOK, []byte(`segments are added/deleted to user`))
}

// @Summary		add segments to user with deadline
// @Description	adding segments to user with expire date (subscription)
// @Tags			segment
// @Param			ttl 	path		int					true	"ttl"
// @Param			segmentName 	path		string					true	"segment name"
// @Param			userId	path		int					true	"User id"
// @Success		200		"Segments for user:"
// @Router			/segment/add_user_deadline [post]
func (s *Server) AddUserDeadline(w http.ResponseWriter, r *http.Request) {
	var input model.DeadlineInput
	data, err := io.ReadAll(r.Body)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`error in reading a body`))
		return
	}

	if err = json.Unmarshal(data, &input); err != nil {
		jsonRespond(w, http.StatusBadRequest, []byte(`error in unmarshalling a body`))
		return
	}
	fmt.Println(&input)
	err = s.Service.AddSegmentDeadline(input.Ttl, input.SegmentName, input.UserId)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, []byte(`cannot add/delete segments to user`))
		return
	}

	jsonRespond(w, http.StatusOK, []byte(`segments with deadline added`))
}
