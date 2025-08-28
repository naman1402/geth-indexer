# Ethereum Event Indexer

A Go-based application that indexes and stores Ethereum smart contract events in real-time. Built to handle both historical and live events, with support for PostgreSQL storage.

## Features

- Real-time event monitoring
- Historical event indexing
- Configurable block range scanning
- PostgreSQL storage integration
- Environment-based configuration
- Docker support

## Architecture

The application consists of three main components:

1. **Subscriber**: Connects to Ethereum nodes, fetches contract ABIs, and listens for events
2. **Indexer**: Processes events and stores them in PostgreSQL
3. **CLI**: Handles configuration and command-line arguments

## Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)
- Ethereum node access (Infura/Local)
- Etherscan API key

## Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/naman1402/geth-indexer.git
cd geth-indexer
```

2. Copy and configure environment variables:
```bash
cp .env.example .env
```

3. Edit `.env` file with your credentials:
```env
# Blockchain Configuration
CONTRACT_ADDRESS=your_contract_address
START_BLOCK=0
END_BLOCK=latest
EVENT_NAME=Transfer

# API Keys
ETHERSCAN_API_KEY=your_etherscan_api_key
INFURA_API_KEY=your_infura_api_key

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=geth_indexer

# RPC Configuration
RPC_URL=wss://mainnet.infura.io/ws/v3/${INFURA_API_KEY}
```

## Running the Application

### Using Docker (Recommended)

1. Build and start the containers:
```bash
docker compose build --no-cache
docker compose up -d
```

2. View logs:
```bash
docker compose logs -f app
```

3. Stop the application:
```bash
docker compose down
```

### Running Locally

1. Install dependencies:
```bash
go mod download
```

2. Build the application:
```bash
go build -o geth-indexer
```

3. Run the indexer:
```bash
./geth-indexer Transfer Approval
```

## Development

### Git Workflow

1. Create a new branch:
```bash
git checkout -b feature/your-feature-name
```

2. Commit changes:
```bash
git add .
git commit -m "feat: your feature description"
```

3. Push changes:
```bash
git push origin feature/your-feature-name
```

### Code Structure

```
geth-indexer/
├── cli/         # Command-line interface and configuration
├── indexer/     # Event processing and storage
├── subscriber/  # Ethereum node connection and event listening
├── .env.example # Environment template
└── docker-compose.yaml
```

## Event Data

The indexer stores the following information for each event:
- Event Name
- Block Number
- Block Hash
- Contract Address
- Event-specific data (e.g., from, to, value for Transfer events)

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

