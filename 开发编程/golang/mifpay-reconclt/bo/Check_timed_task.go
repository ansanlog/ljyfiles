package bo

import "time"

type Check_timed_task struct {
	Id             int
	Name           string
	Create_time    time.Time
	Type           int
	Separate_time  int
	Status         int
	Remark         string
	Func           string
	func_type      int
	Start_time     time.Time
	Next_exec_time time.Time
	Last_exec_time time.Time
	Is_available   int
}
