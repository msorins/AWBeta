package solvers

type ISolver interface {
	updateStatuses()

	GetStatuses() []string
	GetLastStatus() [] string
}