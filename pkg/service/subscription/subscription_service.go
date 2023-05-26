package subscription

import (
	"btc-uah-rates/pkg/repository"
	"errors"
)

type SubscriptionService interface {
	AddSubscription(email string) error
	IsSubscribed(email string) (bool, error)
	GetAll() ([]*string, error)
}

type subscriptionServiceImpl struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionServiceImpl{
		repo: repo,
	}
}

func (s *subscriptionServiceImpl) AddSubscription(email string) error {
	isSubscribed, err := s.IsSubscribed(email)
	if err != nil {
		return err
	}

	if isSubscribed {
		return errors.New("already subscribed")
	}

	return s.repo.Add(email)
}

func (s *subscriptionServiceImpl) IsSubscribed(email string) (bool, error) {
	exists, err := s.repo.Exists(email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *subscriptionServiceImpl) GetAll() ([]*string, error) {
	return s.repo.GetAll()
}
