package utils

import "ClockOut/internal/model"

func NegateStatus(status model.Type) model.Type {
	if status == "start" {
		return "finish"
	}
	return "start"
}
