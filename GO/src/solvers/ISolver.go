package solvers

type ISolver interface {
	GetStatusesForAwb(awb string)

	GetStatuses()
	GetLastStatus()
}