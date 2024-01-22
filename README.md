# "Word of Wisdom" TCP server protected from DDoS attacks with the Proof of Work algorithm

## Task description

* TCP server should be protected from DDOS attacks with the Proof of Work, the challenge-response protocol should be
  used.
* The choice of the POW algorithm should be explained.
* After Proof Of Work verification, server should send one of the quotes from "word of wisdom" book or any other
  collection of the quotes.
* Docker file should be provided both for the server and for the client that solves the POW challenge

## Choice of PoW algorithm

I choose [Hashcash](https://en.wikipedia.org/wiki/Hashcash) algorithm because

1. It's easy to implement
2. It's secure
3. There are a lot of iformation about it
4. Easy to adjust complexity by changing number of leading zeros

## Prerequisites

* Go 1.21+
* Docker

## How to execute

```
# run tests
make run-tests

# start server
make run-server

# start client
make run-client

# start server and client in docker compose
make deploy
```