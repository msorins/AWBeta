package subscription

import (
	"solvers"
	"errors"
	"fmt"
	"time"
)

type SubscriptionManagerEntity struct {
	Solver               solvers.ISolver
	LastNumberOfEntities int
	ByUserId             string
}

type SubscriptionManager struct {
	subscriptions map[string] SubscriptionManagerEntity
}

func SubscriptionManagerBuilder() SubscriptionManager {
	sm := SubscriptionManager{map[string] SubscriptionManagerEntity{}}

	return sm
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

func (manager *SubscriptionManager) CheckForChanges() (map[string] []string, error) {
	fmt.Println("Starting to check for subscriptions - " + time.Now().String())
	sendingmsgs := map[string] []string{}

	for _, subscription := range manager.subscriptions {
		// Must check for update
		oldNrOfStatuses := subscription.LastNumberOfEntities
		statuses, responseCode := subscription.Solver.GetStatuses()

		diff := len(statuses) - oldNrOfStatuses
		// Means the the status has updated => send user the update
 		if responseCode == solvers.SOLVER_OK && diff > 1{
			msgs := []string{"A new update for awb " + subscription.Solver.GetAwb()}
			for _, status := range statuses[ len(statuses) - diff :] {
				msgs = append(msgs, status)
			}
			sendingmsgs[subscription.ByUserId] = msgs

			// Update the last number of entities
			subscription.LastNumberOfEntities = len(statuses)
		}
	}

	return sendingmsgs, nil
}