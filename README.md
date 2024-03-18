# learn-cabin
Play around with Cabin

# Synchronizing Casbin Enforcers with Kafka

This is the idea:
![Casbin Enforcers synchronization using Kafka](assets/casbin-kafka.svg "Casbin Enforcers synchronization using Kafka")

## How to set it up
1. Set up Kafka by running:

```bash
make kafka
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

4. On Shell 1, press enter to simulate an API call that triggers policy update.
5. You should see in Shell 2 that the consumer receives a message to update its policy.
  - I've saved consumer's policy in `./consumer/consumer_policy.csv`. You can compare it with the producer's policy which is `./config/initial_policy.csv`.
