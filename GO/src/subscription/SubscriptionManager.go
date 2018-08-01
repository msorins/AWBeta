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
	subscriptions map[string] SubscriptionManagerEntity // awb -> subscriptionEntity
}

func SubscriptionManagerBuilder() SubscriptionManager {
	sm := SubscriptionManager{map[string] SubscriptionManagerEntity{}}

	return sm
}

func (manager *SubscriptionManager) GetSubscription(awb string) (SubscriptionManagerEntity, error) {
	_, found := manager.subscriptions[awb]
	if found == false {
		return SubscriptionManagerEntity{}, errors.New(fmt.Sprintf("Awb %s not found in the SubscriptionManager", awb))
	}

	return manager.subscriptions[awb], nil
}

func (manager *SubscriptionManager) AddSubscription(awb string, subscription SubscriptionManagerEntity) error {
	manager.subscriptions[awb] = subscription

	return nil
}

func (manager *SubscriptionManager) RemoveSubscription(awb string) error {
	_, found := manager.subscriptions[awb]
	if found == false {
		return errors.New(fmt.Sprintf("Awb %s not found in the SubscriptionManager", awb))
	}

	delete(manager.subscriptions, awb)

	return nil
}

func (manager *SubscriptionManager) updateSubscription(awb string, newSubscription SubscriptionManagerEntity) error {
	_, found := manager.subscriptions[awb]
	if found == false {
		return errors.New(fmt.Sprintf("Awb %s not found in the SubscriptionManager", awb))
	}

	manager.subscriptions[awb] = newSubscription

	return nil
}

func (manager *SubscriptionManager) IdExists(awb string) bool {
	_, found := manager.subscriptions[awb]

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