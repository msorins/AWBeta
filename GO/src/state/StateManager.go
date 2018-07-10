package state

import (
	"errors"
	"fmt"
	"solvers"
)

type StateManagerStruct struct {
	Solver solvers.ISolver
	State  UserState
}

type StateManager struct {
	states map[string] StateManagerStruct
}

func StateManagerBuilder() StateManager {
	st := StateManager{states: make(map[string] StateManagerStruct)}

	return st
}

func (stateObj *StateManager) SetState(id string, solver solvers.ISolver, state UserState) error {
	var newState = StateManagerStruct{Solver:solver, State:state}
	stateObj.states[id] = newState

	return nil
}


func (stateObj *StateManager) GetState(id string) (StateManagerStruct, error) {
	_, found := stateObj.states[id]
	if found == false {
		return StateManagerStruct{}, errors.New(fmt.Sprintf("Id %s not found in the StateManager", id))
	}

	return stateObj.states[id], nil
}

func (stateObj *StateManager) UpdateState(id string, newState UserState) error {
	_, found := stateObj.states[id]
	if found == false {
		return errors.New(fmt.Sprintf("Id %s not found in the StateManager", id))
	}

	var stateOfId = stateObj.states[id]
	stateOfId.State = newState
	stateObj.states[id] = stateOfId

	return nil
}

func (stateObj *StateManager) RemoveState(id string) error {
	_, found := stateObj.states[id]
	if found == false {
		return errors.New(fmt.Sprintf("Id %s not found in the StateManager", id))
	}

	delete(stateObj.states, id)

	return nil
}

func (stateObj *StateManager) UpdateSolver(id string, newSolver solvers.ISolver) error {
	_, found := stateObj.states[id]
	if found == false {
		return errors.New(fmt.Sprintf("Id %s not found in the StateManager", id))
	}

	var stateOfId = stateObj.states[id]
	stateOfId.Solver = newSolver
	stateObj.states[id] = stateOfId

	return nil
}

func (stateObj *StateManager) IdExists(id string) bool {
	_, found := stateObj.states[id]
	if found == false {
		return false;
	}

	return true
}