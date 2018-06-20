package solvers

type ISolver interface {
	GetStatusesForAwb() []IPackageStatus

	GetStatuses() []IPackageStatus
	GetLastStatus() IPackageStatus
}