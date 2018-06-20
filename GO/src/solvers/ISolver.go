package solvers

type ISolver interface {
	getStatuses(awb string)
	getLastStatus(awb string)
}