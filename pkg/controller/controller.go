package controller

import (
	"fmt"
	"log"
	"net/http"

	"btc-uah-rates/pkg/service/email"
	"btc-uah-rates/pkg/service/notification"
	"btc-uah-rates/pkg/service/rate"
	"btc-uah-rates/pkg/service/subscription"
	"btc-uah-rates/pkg/utils"
)

type Controller struct {
	RateService         rate.RateService
	EmailService        email.EmailService
	SubscriptionService subscription.SubscriptionService
	NotificationService notification.NotificationService
}

func NewController(rateService rate.RateService, emailService email.EmailService,
	subscriptionService subscription.SubscriptionService, notificationService notification.NotificationService) *Controller {
	return &Controller{
		RateService:         rateService,
		EmailService:        emailService,
		SubscriptionService: subscriptionService,
		NotificationService: notificationService,
	}
}

func (c *Controller) GetRateHandler(w http.ResponseWriter, r *http.Request) {
	currentRate, done := c.processRate(w)
	if done {
		return
	}

	c.prepareRateResponse(w, currentRate)
}

func (c *Controller) SubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	providedEmail, done := utils.ExtractFormValue(w, r, "email")
	if done {
		return
	}

	if c.processNewSubscription(w, providedEmail) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) SendEmailsHandler(w http.ResponseWriter, r *http.Request) {
	sent := c.NotificationService.Notify()
	if sent {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *Controller) prepareRateResponse(w http.ResponseWriter, currentRate float64) {
	jsonResponse := []byte(fmt.Sprintf("%.2f", currentRate))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (c *Controller) processRate(w http.ResponseWriter) (float64, bool) {
	currentRate, err := c.RateService.GetRate()
	if err != nil {
		log.Println("Error getting rate:", err)
		http.Error(w, "Failed to get rate", http.StatusInternalServerError)
		return 0, true
	}
	return currentRate, false
}

func (c *Controller) processNewSubscription(w http.ResponseWriter, providedEmail string) bool {
	err := c.SubscriptionService.AddSubscription(providedEmail)
	if err != nil {
		if err.Error() == "already subscribed" {
			log.Println(providedEmail, err)
			http.Error(w, "", http.StatusConflict)
		} else {
			log.Println("Error adding subscription:", err)
			http.Error(w, "Failed to add subscription", http.StatusInternalServerError)
		}
		return true
	}
	return false
}
