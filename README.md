## 📜 Word Of Wisdom

A TCP server implementation that sends a quote from the "Word of Wisdom" book after verifying Proof of Work (PoW) solved challenge to prevent DDOS attacks.


## 🔒 PoW Implementation

The server uses a PoW algorithm to protect itself from potential DDOS attacks. The challenge-response protocol requires the client to provide a solution that satisfies the PoW algorithm in order to receive a quote from the server. The algorithm is chosen based on HMAC-SHA256 hash function.


## 💬 Sending Quotes

After successfully verifying the PoW challenge, the server sends a quote from the "Word of Wisdom" book. The quotes can be easily modified and added based on your desired collection. For now they are hardcoded in config file but storing them in separate database is in roadmap.

## 🐳 Dockerization

Both the server and client that solves the PoW challenge are designed to be Dockerized environments. This ensures the ease of deployment, consistency, and portability of the system.

## 🚀 Getting Started

To start using the TCP server, follow the instructions below:

1) Clone this repository
```bash
$ git clone https://github.com/andrew528i/wisdom-server.git
```

2) Build the Docker images for the server and the client
```bash
$ cd docker
$ docker-compose build server
$ docker-compose build client
```

3) Run the TCP server
```bash
$ docker-compose up server
```

4) Run the PoW client to solve the challenge and get a quote
```bash
$ docker-compose run --rm client
```

## 🧪 Running unit tests
```bash
$ go test -v ./...
=== RUN   TestChallenge_Solve_Check
=== RUN   TestChallenge_Solve_Check/all_is_ok
=== RUN   TestChallenge_Solve_Check/swap_secret
=== RUN   TestChallenge_Solve_Check/exceed_nonce
=== RUN   TestChallenge_Solve_Check/increase_difficulty
=== RUN   TestChallenge_Solve_Check/make_solution_invalid
=== RUN   TestChallenge_Solve_Check/exceed_signature_deadline
=== RUN   TestChallenge_Solve_Check/spoof_data
=== RUN   TestChallenge_Solve_Check/spoof_nonce
=== RUN   TestChallenge_Solve_Check/spoof_deadline
--- PASS: TestChallenge_Solve_Check (2.47s)
    --- PASS: TestChallenge_Solve_Check/all_is_ok (0.00s)
    --- PASS: TestChallenge_Solve_Check/swap_secret (0.00s)
    --- PASS: TestChallenge_Solve_Check/exceed_nonce (0.00s)
    --- PASS: TestChallenge_Solve_Check/increase_difficulty (0.00s)
    --- PASS: TestChallenge_Solve_Check/make_solution_invalid (0.00s)
    --- PASS: TestChallenge_Solve_Check/exceed_signature_deadline (2.00s)
    --- PASS: TestChallenge_Solve_Check/spoof_data (0.00s)
    --- PASS: TestChallenge_Solve_Check/spoof_nonce (0.00s)
    --- PASS: TestChallenge_Solve_Check/spoof_deadline (0.00s)
PASS
ok  	andrew528i/wisdom_server/internal	3.031s
```

## 📝 Note

This is a sample implementation of a TCP server protected with PoW challenge-response protocol. Please use caution when deploying to production and make necessary modifications based on your specific requirements.