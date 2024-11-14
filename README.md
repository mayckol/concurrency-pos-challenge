# Auction Application

## Overview

This project is an auction application where users can create auctions and place bids. The main objective of this task was to add a new feature that automatically closes auctions after a predefined time using Go routines (goroutines).

## What Was Done

To complete the test, the following steps were taken:

1. **Implemented a function to calculate the auction duration based on environment variables.**
    - Added a function `getAuctionDuration()` to retrieve the auction duration from the `AUCTION_DURATION` environment variable.
    - Set a default duration in case the environment variable is not defined.

2. **Created a new goroutine to automatically close expired auctions.**
    - Implemented a goroutine within the `startAuctionCloser()` function that periodically checks for expired auctions and updates their status to `Completed`.
    - Used a ticker to run the check at regular intervals, ensuring efficient resource usage.

3. **Developed a test to validate the automatic closure of auctions.**
    - Wrote a test `TestAuctionAutoClose` to create an auction with a short duration, wait for it to expire, and verify that its status is updated accordingly.

## How to Run the Application

### Prerequisites

- **Docker** and **Docker Compose** installed on your system.

### Environment Variables

Create a `.env` file in the `cmd/auction` directory with the following content:

```env
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=20s

MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://localhost:27017
MONGODB_DB=auctions
AUCTION_INTERVAL=20s
AUCTION_DURATION=5m

```

### Running the Application

1. Build and Start the Services: In the root directory of the project, run:

```bash
docker-compose up --build
```

2. Access the Application: Open your browser and go to `http://localhost:8080`.

## Running the Tests

To run the tests, execute the following command:

```bash
docker-compose run app go test ./...
```

This will execute all the tests in the application, including the test for automatic auction closure.

Notes
MongoDB Service: 
- Ensure that the MongoDB service (mongodb_challenge) is up and running before starting the application.
- Environment Variables: The application reads configuration from environment variables defined in the .env file.
- Ports: If the default ports (8080 for the application and 27017 for MongoDB) are in use, you may need to adjust them in the docker-compose.yml file.