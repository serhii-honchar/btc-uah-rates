package notification

import (
	"btc-uah-rates/pkg/service/email"
	"btc-uah-rates/pkg/service/rate"
	"btc-uah-rates/pkg/service/subscription"
	"fmt"
	"log"
	"time"
)

type NotificationService interface {
	Notify() bool
}

type notificationServiceImpl struct {
	subscriptionService subscription.SubscriptionService
	emailService        email.EmailService
	rateService         rate.RateService
}

func NewNotificationService(subscriptionService subscription.SubscriptionService,
	emailService email.EmailService, rateService rate.RateService) NotificationService {
	return &notificationServiceImpl{
		subscriptionService: subscriptionService,
		emailService:        emailService,
		rateService:         rateService,
	}
}

func (n *notificationServiceImpl) Notify() bool {
	recipients := n.getRecipients()

	currentRate, err := n.rateService.GetRate()
	if err != nil {
		log.Println("Error while getting current rate")
	}

	messageWithTime := n.prepareSubject()
	message := n.prepareMessage(currentRate)
	err = n.emailService.SendEmails(messageWithTime, message, recipients)
	if err != nil {
		return false
	}
	return true
}

func (n *notificationServiceImpl) getRecipients() []string {
	subscriptionsPtr, err := n.subscriptionService.GetAll()
	if err != nil {
		log.Println("Error while getting subscriptions")
	}
	subscriptions := make([]string, len(subscriptionsPtr))
	for i, subPtr := range subscriptionsPtr {
		subscriptions[i] = *subPtr
	}
	return subscriptions
}

func (n *notificationServiceImpl) prepareSubject() string {
	formattedTime := time.Now().Format("02-01-2006 15:04:05")
	return fmt.Sprintf("%s - %s", "BTC/UAH rate", formattedTime)
}

func (n *notificationServiceImpl) prepareMessage(currentRate float64) string {
	return fmt.Sprintf("BTC/UAH rate:%.2f \n\n\nYou received this message because you're lucky", currentRate)
}
