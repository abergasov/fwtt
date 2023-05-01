## Overview
### Description
Design and implement “Word of Wisdom” tcp server.  
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
• The choice of the POW algorithm should be explained.  
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
• Docker file should be provided both for the server and for the client that solves the POW challenge

## local run 
```shell
# run docker container with server
make run
```

run js client which solves the POW challenge
```shell
cd client && yarn install && cd ..
make client
```

tests and linter 
```shell
make test
make lint
make bench
```

## Solution
The main idea is that the server generates a random challenge and sends it to the client. The challenge is represented as a JSON object that includes an array of one or more challenge strings, a difficulty value, and an algorithm used for hashing.

If the algorithm is sha256, the JSON object will look like this:
```json
{
  "challenges": [
    "74745c0d-0274-4f61-a1fd-96088c6044fe",
    "0336b71f-d6ca-4b73-a1c2-eb7c9ae5fff2"
  ],
  "difficulty": 2,
  "algorithm": "sha256"
}
```

If the algorithm is scrypt, the JSON object will look like this:

```json
{
  "challenges": [
    "ea8abc2c-51f3-4dec-97e2-b94a224e1e34"
  ],
  "difficulty": 2,
  "algorithm": "scrypt",
  "algo_params": {
    "n": 1024,
    "r": 1,
    "p": 1,
    "key_len": 32
  }
}
```
The endpoint that returns the random challenge can be hidden behind a rate limiter to prevent abuse.

The client solves the challenge and sends the solution back to the server. The difficulty of the solution can be a floating-point value that depends on the current state of the server. 

For example, if the server is consuming a low amount of resources, the difficulty can be low. If the server is under heavy load, the difficulty can be increased. If the server is under super heavy load, the POW algorithm can be replaced with a more heavy one, for example, switching from sha256 to scrypt.

### POW algorithm
* SHA256 - is a widely-used and fast hash function, suitable for use when the server is under heavy load. It requires fewer resources than Scrypt, making it a good option when the server is under pressure. The difficulty level can be easily increased to adjust to changing server conditions.
* SCRYPT - is a memory-hard hash function that can be used when server load is low. It is ideal when server load is below a certain threshold, and difficulty level can be easily increased to adapt to changing server conditions. However, it requires more memory than SHA256.
### Improvements
* The current implementation is suitable for a single server, but for a distributed environment with multiple servers, a distributed proof-of-work (POW) algorithm needs to be implemented to ensure the system's security and stability. This can be achieved by distributing the computation workload among multiple servers and coordinating the computation results to verify the correctness of the computation. 

