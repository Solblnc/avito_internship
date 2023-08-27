package model

type User struct {
	Id       int
	Segments []Segment
}

type Segment struct {
	Id   int
	Name string
}
