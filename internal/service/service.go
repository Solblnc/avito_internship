package service

import (
	"errors"
	"fmt"
)

var (
	ErrCreateSegment = errors.New("error in creating a segment")
	ErrDeleteSegment = errors.New("error in deleting a segment")
	ErrAddUser       = errors.New("error in addind segments to user")
	ErrCreateUsers   = errors.New("error in creating users")
	ErrGetHistory    = errors.New("error in getting a history (csv file)")
)

type Store interface {
	Create(segment string, percent uint) (int, error)
	Delete(segment string) error
	AddUser(segmentsAdd []string, segmentsDelete []string, userId int) error
	AddSegmentDeadline(ttl int, segmentName string, userId int) error
	GetActualSegments(userId int) []string
	CreateUser() error
	GetHistory(year, month int) (string, error)
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{Store: store}
}

func (s *Service) Create(segment string, percent uint) (int, error) {
	id, err := s.Store.Create(segment, percent)
	if err != nil {
		fmt.Println(err)
		return 0, ErrCreateSegment
	}
	return id, nil
}

func (s *Service) Delete(segment string) error {
	if err := s.Store.Delete(segment); err != nil {
		fmt.Println(err)
		return ErrDeleteSegment
	}
	return nil
}

func (s *Service) AddUser(segmentsAdd []string, segmentsDelete []string, userId int) error {
	if err := s.Store.AddUser(segmentsAdd, segmentsDelete, userId); err != nil {
		fmt.Println(err)
		return ErrAddUser
	}
	return nil
}

func (s *Service) GetActualSegments(userId int) []string {
	segments := s.Store.GetActualSegments(userId)
	return segments
}

func (s *Service) CreateUser() error {
	return s.Store.CreateUser()
}

func (s *Service) GetHistory(year, month int) (string, error) {
	history, err := s.Store.GetHistory(year, month)
	if err != nil {
		fmt.Println(err)
		return "", ErrGetHistory
	}

	return history, nil
}

func (s *Service) AddSegmentDeadline(ttl int, segmentName string, userId int) error {
	return s.Store.AddSegmentDeadline(ttl, segmentName, userId)
}
