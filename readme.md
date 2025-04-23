<!-- PROJECT LOGO -->
<br />
<div align="center">
  <h3 align="center">Qiscus Custom Agent Allocation</h3>

  <p align="center">
    Custom chat allocation to agent for qiscus service
    <br />
</p>
</div>

<!-- ABOUT THE PROJECT -->

## About The Project

This is a service to allocate chat to the agent available, eliminating the process of manual assigning to the agent

Currently it can only provide the same amount of max customer of all agent. It can be modified from redis and the default amount of it in the env file

### Built With

1. Golang - Handle most of the process
2. Redis - Queue system of the webhook
3. Docker - Providing self hosted Redis database

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

How to setup the service

### Prerequisites

- Golang  
  Can be installed from this link https://go.dev/doc/install
- Docker
  Can refer to this link https://docs.docker.com/get-started/get-docker/

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/victor-divo/custom-agent-service.git
   ```
2. Install go packages
   ```sh
   go mod tidy
   ```
3. Copy .env.example to `.env`
   ```sh
   cp .env.example .env
   ```
4. Run the redis database using docker or make command
   ```sh
   docker compose up -d
   ```
   or
   ```sh
   make up
   ```
5. Open the redis browser http://localhost:5540 and add max agent config.

   Add connection to Redis database (default url redis://default@redis:6379)

   In redis cli change n to desired number

   ```cli
   SET config:max_agent_chat n
   ```

6. Run the service
   ```sh
   go run cmd/main.go
   ```
   If you want to debug, run the default vscode debuger

<!-- USAGE EXAMPLES -->

## How It Works?

### Flow Chart

1. WebHook Receiver.  
   This is the one who tell the system if any chat without an Agent is coming, and add it to the queue  
   ![Flow Chart](https://raw.githubusercontent.com/victor-divo/custom-agent-service/main/documentation/diagram_webhook.png)
2. Worker.  
   This is the code that process all the item in the queue
   ![Flow Chart](https://raw.githubusercontent.com/victor-divo/custom-agent-service/main/documentation/diagram_worker.png)

### Sequence Diagram

This is the sequence diagram of how the general process happens continously
![Sequence Diagram](https://raw.githubusercontent.com/victor-divo/custom-agent-service/main/documentation/sequence_diagram.png)

<!-- CONTACT -->

## Contact

Victor Divo Mahendra - victormahendrausm@gmail.com - [Whatsapp](wa.me/+6287776901628)
