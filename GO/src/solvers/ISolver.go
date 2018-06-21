package solvers

type ISolver interface {
	GetStatusesForAwb() []IPackageStatus

	GetStatuses() []string
	GetLastStatus() [] string
}