package model

type User struct {
	Id int
}

type UserSegments struct {
	UserId    int `json:"user_id"`
	SegmentId int `json:"segment_id"`
}

type Segment struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Percent uint   `json:"percent"`
}

type InputAddUser struct {
	SegmentsAdd    []string `json:"segments_add"`
	SegmentsDelete []string `json:"segments_delete"`
	UserId         int      `json:"user_id"`
}

type DeadlineInput struct {
	Ttl         int    `json:"ttl"`
	SegmentName string `json:"segment_name"`
	UserId      int    `json:"user_id"`
}
