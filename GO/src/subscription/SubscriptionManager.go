package subscription

import (
	"solvers"
	"time"
	"errors"
)

type SubscriptionManagerEntity struct {
	Solver solvers.ISolver
	lastNumberOfEntities int
	lastUpdateCheck	time.Time
}

type SubscriptionManager struct {
	subscriptions map[string] SubscriptionManagerEntity
}

func (manager *SubscriptionManager) GetSubscription(id string, subscription SubscriptionManagerEntity) (SubscriptionManagerEntity, error) {
	return SubscriptionManagerEntity{}, errors.New("Not yet implemented")
}

func (manager *SubscriptionManager) AddSubscription(id string, subscription SubscriptionManagerEntity) error {
	return errors.New("Not yet implemented")
}

func (manager *SubscriptionManagerEntity) RemoveSubscription(id string) error {
	return errors.New("Not yet implemented")
}

func (manager *SubscriptionManagerEntity) updateSubscription(id string, newSubscription SubscriptionManagerEntity) error {
	return errors.New("Not yet implemented")
}

func (manager *SubscriptionManagerEntity) IdExists(id string) bool {
	return false
}

func (manager *SubscriptionManagerEntity) CheckForChanges() error {
	return errors.New("Not yet implemented")
}