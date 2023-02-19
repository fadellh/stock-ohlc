# Stockbit Test

## System Design
![high-level-design](./ohlc-high-level-design.png)

- **Request-genator** : For simplicity in this project. Let's think the request-generator service as something that can handle a lot of transactions with tiny, speedy data. We'll be getting a lot of requests, then we'll send them on to the message broker.
- **Calculation Service** : We've got an amazing calculator service that'll make your calculations so much easier! 🤩
- **Summary Service** : Where clients can get instant results!

