# learn-casbin
Play around with Casbin

# Synchronizing Casbin Enforcers with RabbitMQ

This is the idea:
![Casbin Enforcers synchronization using RabbitMQ](assets/casbin-rabbitmq.svg "Casbin Enforcers synchronization using RabbitMQ")

## How to set it up
1. Set up RabbitMQ by running:

```bash
make rabbitmq
```

2. Set up producer by running:

```bash
#Shell 1
make producer
```

3. Open a new shell. Set up consumer by running:

```bash
#Shell 2
make consumer
```

4. In Shell 1, press enter to simulate an API call that triggers policy update.
5. You should see in Shell 2 that the consumer receives a message to update its policy.
  - I've made consumer save its policy into `./consumer/consumer_policy.csv`. You can compare it with the producer's policy which is `./config/initial_policy.csv`.
