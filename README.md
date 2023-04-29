## Overview
### Description
Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.  
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
• The choice of the POW algorithm should be explained.  
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
• Docker file should be provided both for the server and for the client that solves the POW challenge

## local run 
```shell
make run
```

run js client which solves the POW challenge
```shell
make client
```

tests and linter 
```shell
make test
make lint
make bench
```

## Solution
### Description
### POW algorithm
* SHA256
* SCRYPT
### Improvements
