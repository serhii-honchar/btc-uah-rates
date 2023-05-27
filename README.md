# btc-uah-rates
Application perform the next actions:
1. Gets current rate of BTC/UAH (**GET /api/rate**)
Public API https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=UAH is used by default and can be overriden by RATE_PROVIDER_URL environment variable
2. Stores list of subsribed user (**POST /api/subscribe**, email should be send in form)
Storage implemented using filesystem, each subscription is a new line in text file. By default - subscriptions.txt file is used in project root directory (can be overriden by DATABASE_URL property)
The concurrent access to the file is organized using Mutual exclusion lock.
The duplicated subscriptions isn't allowed and 409 error is returned when user tries to add the same email again (so far the entered value isn't trimmed)
3. Sends emails to the list of stored subscribers with current BTC/UAH rate (**POST /api/sendEmails**)
Email notifications are implemented using Mailhog SMTP server.
SMTP connection parameters are configured using environment variables that can be overriden (EMAIL_HOST, EMAIL_PORT, EMAIL_USERNAME, EMAIL_PASSWORD)


There's a Dockerfile in project root directory that can be used to build and start application:

`docker build . -t rate-tracker`

`docker run -d -p 9090:8080 rate-tracker`


Also, the easiest way to start application for testing is to use docker compose:

`/usr/bin/docker compose -f docker-compose.yaml up -d`

It will prepare and start 2 services: 
1) application (will be accessable on localhost:9090)
2) Mailhog SMTP server (UI will be available using localhost:8025 externally or mailhog:8025 within docker-compose) 

Example of received email:
![image](https://github.com/serhii-honchar/btc-uah-rates/assets/8502717/4759e976-e3e0-4ca2-ad69-38ba69ebc932)
