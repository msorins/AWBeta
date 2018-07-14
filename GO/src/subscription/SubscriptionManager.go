package subscription

import (
	"solvers"
	"errors"
	"fmt"
)

type SubscriptionManagerEntity struct {
	Solver solvers.ISolver
	lastNumberOfEntities int
}

type SubscriptionManager struct {
	subscriptions map[string] SubscriptionManagerEntity
}

func (manager *SubscriptionManager) GetSubscription(id string) (SubscriptionManagerEntity, error) {
	_, found := manager.subscriptions[id]
	if found == false {
		return SubscriptionManagerEntity{}, errors.New(fmt.Sprintf("Id %s not found in the SubscriptionManager", id))
	}

	return manager.subscriptions[id], nil
}

func (manager *SubscriptionManager) AddSubscription(id string, subscription SubscriptionManagerEntity) error {
	manager.subscriptions[id] = subscription

	return nil
}

func (manager *SubscriptionManager) RemoveSubscription(id string) error {
	_, found := manager.subscriptions[id]
	if found == false {
		return errors.New(fmt.Sprintf("Id %s not found in the SubscriptionManager", id))
	}

	delete(manager.subscriptions, id)

	return nil
}

func (manager *SubscriptionManager) updateSubscription(id string, newSubscription SubscriptionManagerEntity) error {
	_, found := manager.subscriptions[id]
	if found == false {
		return errors.New(fmt.Sprintf("Id %s not found in the SubscriptionManager", id))
	}

	manager.subscriptions[id] = newSubscription

	return nil
}

func (manager *SubscriptionManager) IdExists(id string) bool {
	_, found := manager.subscriptions[id]

	return found
}

func (manager *SubscriptionManager) CheckForChanges() error {
	for _, subscription := range manager.subscriptions {
		// Must check for update
		oldNrOfStatuses := subscription.lastNumberOfEntities
		statuses, responseCode := subscription.Solver.GetStatuses()

		// Means the the status has updated => send user the update
		diff := len(statuses) - oldNrOfStatuses
 		if responseCode == solvers.SOLVER_OK && diff > 1{
			// TO DO -> send user the new statuses

			//
			subscription.lastNumberOfEntities = len(statuses)
		}
	}

	return nil
}