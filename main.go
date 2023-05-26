package main

import (
	"btc-uah-rates/pkg/controller"
	"btc-uah-rates/pkg/repository"
	"btc-uah-rates/pkg/service/email"
	"btc-uah-rates/pkg/service/notification"
	"btc-uah-rates/pkg/service/rate"
	"btc-uah-rates/pkg/service/subscription"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	fileRepository := repository.NewFileRepository()
	rateService := rate.NewRateService()
	subscriptionService := subscription.NewSubscriptionService(fileRepository)
	emailService := email.NewEmailService()
	notificationService := notification.NewNotificationService(subscriptionService, emailService, rateService)
	rateController := controller.NewController(rateService, emailService, subscriptionService, notificationService)

	router := mux.NewRouter()
	router.HandleFunc("/rate", rateController.GetRateHandler).Methods(http.MethodGet)
	router.HandleFunc("/subscribe", rateController.SubscriptionHandler).Methods(http.MethodPost)
	router.HandleFunc("/sendEmails", rateController.SendEmailsHandler).Methods(http.MethodPost)

	http.Handle("/api/", http.StripPrefix("/api", router))

	log.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
