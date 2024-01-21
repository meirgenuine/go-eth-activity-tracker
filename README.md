# go-eth-activity-tracker
Activity Tracker for Ethereum Mainnet

## Overview
`go-eth-activity-tracker` is a Go (Golang) application designed to track and calculate activity metrics for Ethereum addresses on the Ethereum Mainnet. It provides insights into the most active Ethereum addresses based on the number of ERC20 token transfers (both incoming and outgoing).

*__Note:__ <u>Activity count</u> for the specific address is incremented if an incoming or outgoing ERC20 token transfer occurs. This means that if an address has no ERC20 token transfers, it will not be included in the results.*

## Features
- Retrieves and analyzes data from the Ethereum Mainnet.
- Identifies the most active Ethereum addresses based on configured criteria.
- Rate limits requests to ensure efficient and responsible data retrieval.

## Setup

Before running the application, follow these steps to set up your environment:

### 1. Obtain an Access Token

To access Ethereum data from the Mainnet, you will need an Access Token from GetBlock.io. If you don't have one, sign up on their website to obtain an Access Token.

### 2. Set the Environment Variable

Once you have your Access Token, set it as an environment variable in your terminal. Replace `"your_token_here"` with your actual Access Token:

```bash
export ETH_ACCESS_TOKEN="your_token_here"
```

## Usage

### Running the Application

To run the go-eth-activity-tracker application, use the following command:

```bash
go run cmd/main.go
```

Alternatively, you can build the application and then run the executable:

```bash
go build -o eth-activity-tracker cmd/main.go
./eth-activity-tracker
```

The application will retrieve the latest Ethereum blocks, calculate activity metrics based on ERC20 token transactions, and display the top 5 active addresses based on the configured criteria.

## License

This project is licensed under the MIT License - see the LICENSE file for details.