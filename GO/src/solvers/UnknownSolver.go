package solvers

import "wit"

type unknownCourierSolver struct {
	awb string
	url string
	Statuses []AWbFanCourierCheckpoint
}


func UnknownFanCourierSolverBuilder(awb string, entities map[string][]wit.WitEntity) ISolver{
	awbSolver := unknownCourierSolver{}
	awbSolver.awb = awb
	return &awbSolver
}

func (solver *unknownCourierSolver) updateStatuses()  {


}

func (awbsolver *unknownCourierSolver) GetStatuses() []string {
	results := []string{}

	results = append(results, "Could not found and AWB and link it with a courier company")
	return results
}

func (awbsolver *unknownCourierSolver) GetLastStatus() []string {
	results := []string{}

	results = append(results, "Could not found and AWB and link it with a courier company")
	return results
}


