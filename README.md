# SnappChat

SnappChat is a CLI-based live chat application built with love by Momoein.

## Features
- Real-time communication between users using NATS for messaging.
- Easy to deploy with Docker.

## Prerequisites

Before setting up the project, ensure you have the following installed on your system:

- **Go (Golang)**: Version 1.23 or higher
- **Docker**: For containerized deployment
- **Nats-server**: If you are not using Docker

## Installation and Setup

Follow these steps to set up and run SnappChat:

### Step 1: Clone the Repository
```bash
git clone https://github.com/momoein/snapp-chat.git
cd snapp-chat
```

### Step 2: Configure the Environment

Edit the `server_config.json` and `client_config.json` files to set up the server and client configurations.

server config Example:
```json
{
    "server": {
        "httpPort": 8080
    },
    "nats": {
        "host": "nats://localhost",
        "port": 4222
    }
}
```

client config Example:
```json
{
    "websocketAddr": {
        "host": "localhost",
        "port": 8080,
        "path": "api/v1/ws"
    },
    "httpAddr": {
        "host": "localhost",
        "port": 8080,
        "path": "api/v1/ws/users"
    }
}
```

### Step 3: Install Dependencies

Run the following command to install dependencies:
```bash
go mod tidy
```

### Step 4: Build the Application

To build the server and client executables:
```bash
go build -o bin/server ./cmd/server

go build -o bin/client ./cmd/client
```

### Step 5: Run the Application

#### Option 1: Run Locally

Start the server:
```bash
./bin/server
```

Run the client:
```bash
./bin/client
```

#### Option 2: Use Docker Compose

Run the following command to start the server application with Docker:
```bash
docker-compose up --build
```

### Step 6: Connect and Chat

Once the server is running, you can connect multiple clients and start chatting in real-time.

## Project Structure

The project is organized as follows:

```
.
├── api               # API definitions and handlers
├── app               # Application-level logic
├── build             # Build scripts and Dockerfiles
├── cmd               # Command-line executables for server 
├── config            # Configuration files
├── internal          # Internal libraries and utilities
└── README.md         # Project documentation
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request for any features or fixes you'd like to contribute.

## Contact

For questions or support, please contact the developer at: [momoein711@gmail.com](mailto:momoein711@gmail.com).

