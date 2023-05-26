package repository

import (
	"btc-uah-rates/pkg/utils"
	"bufio"
	"errors"
	"os"
	"sync"
)

type SubscriptionRepository interface {
	GetAll() ([]*string, error)
	Add(email string) error
	Exists(email string) (bool, error)
}

type fileRepository struct {
	filename string
	mu       sync.Mutex
}

func NewFileRepository() SubscriptionRepository {
	filename := utils.GetEnvOrDefault("DATABASE_URL", "subscriptions.txt")
	return &fileRepository{
		filename: filename,
	}
}

func (r *fileRepository) GetAll() ([]*string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var subscriptions []*string
	for scanner.Scan() {
		subscription := scanner.Text()
		subscriptions = append(subscriptions, &subscription)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *fileRepository) Add(subscription string) error {
	exists, err := r.Exists(subscription)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	file, err := os.OpenFile(r.filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(subscription + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (r *fileRepository) Exists(email string) (bool, error) {
	subscriptions, err := r.GetAll()
	if err != nil {
		return false, err
	}

	for _, subscription := range subscriptions {
		if *subscription == email {
			return true, nil
		}
	}

	return false, nil
}
