package model

import "time"

type Type string

const (
	Start  Type = "start"
	Finish Type = "finish"
)

type ClockIn struct {
	ID         int64
	EmployeeId int64
	Timestamp  time.Time
	Type       Type
}
