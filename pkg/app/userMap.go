package app

import (
	"fmt"
	"time"
)

func (a App) LaunchTimer(userID int64) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user, ok := a.users[userID]
	if !ok {
		return fmt.Errorf("id - %d; error - %w", userID, ErrUserDoesNotExist)
	}

	go user.Run()
	return nil
}

func (a App) StopTimer(userID int64) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	_, ok := a.users[userID]
	if !ok {
		return fmt.Errorf("id - %d; error - %w", userID, ErrUserDoesNotExist)
	}

	a.users[userID].Stop()
	delete(a.users, userID)
	return nil
}

func (a App) AddNewUser(userID int64, newWorker worker) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	_, ok := a.users[userID]
	if ok {
		return fmt.Errorf("id - %d; error - %w", userID, ErrUserAlreadyExists)
	}

	a.users[userID] = newWorker
	return nil
}

func (a App) SetupChillingDuration(userID int64, chillingDuration time.Duration) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user, ok := a.users[userID]
	if !ok {
		return fmt.Errorf("id - %d; error - %w", userID, ErrUserDoesNotExist)
	}

	user.chillingPeriod = chillingDuration
	a.users[userID] = user
	return nil
}

func (a App) SetupHustlingDuration(userID int64, hustlingTimePeriod time.Duration) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	user, ok := a.users[userID]
	if !ok {
		return fmt.Errorf("id - %d; error - %w", userID, ErrUserDoesNotExist)
	}

	user.hustlingPeriod = hustlingTimePeriod
	a.users[userID] = user
	return nil
}
