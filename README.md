# Pooled Funds Financial Application

This project implements a microservice architecture for a financial application that allows users to create and manage pooled funds. It provides features such as user registration and authentication, money pool creation, P2P transactions, basic financial reporting, customizable pool permissions, a notification system, and compliance and security measures.

## Architecture

The application follows a microservice architecture with the following key components:

![Microservice Architecture Diagram](https://i.imgur.com/6klxx29.png)

- API Gateway: Built using Golang and the Gin framework, the API Gateway acts as the entry point for client requests and routes them to the appropriate backend services.
- Authentication Service: Utilizes Firebase for secure user registration and authentication.
- Pools & Transactions Service: Handles the creation and management of money pools and facilitates P2P transactions within the platform.
- User & Compliance Service: Implements KYC procedures and manages user information and compliance-related data.
- Galileo API: Integrates with the Galileo FT payment processor for handling financial transactions.
- Plaid API: Integrates with Plaid for AML and KYC compliance checks.

The microservices communicate with each other using gRPC and are backed by a PostgreSQL database.

## Features

The application offers the following features:

- Secure user registration and authentication
- Money pool creation and management
- P2P transactions within the platform
- Basic financial reporting and a financial dashboard
- Customizable pool permissions
- Notification system for important financial activities
- Compliance and security measures, including KYC procedures

## User Flows

![Creating a Digital Wallet Flow](https://i.imgur.com/YjlYtZp.png)

### Creating a Digital Wallet

The user flow for creating a digital wallet involves the following steps:

1. The user initiates the process by providing their personal account information.
2. The application performs identity verification and KYC checks.
3. If the checks pass, the user can log in and retrieve their digital wallet information.
4. If the checks fail, the user is prompted to provide additional information or the process is terminated.

### Creating a Money Pool

The user flow for creating a money pool on an existing user account with a digital wallet includes:

1. The user initiates the process by providing the necessary information for the money pool.
2. The application performs secondary account verification and sends the information to the payment processor for handling.
3. If successful, the money pool is created, and the user can log in and retrieve the money pool information.
4. If unsuccessful, the user is prompted to provide additional information or the process is terminated.

## Technologies Used

- Backend: Golang, Gin Framework, GORM
- Authentication: Firebase
- Payment Processor: Galileo FT
- AML Compliance: Plaid
- KYC Compliance: Plaid
- Database: PostgreSQL
- API Documentation: Swagger/OpenAPI

## Getting Started

To run the application locally, follow these steps:

1. Clone the repository.
2. Install the necessary dependencies.
3. Configure the required environment variables.
4. Start the individual microservices.
5. Access the API Gateway at `http://localhost:3000`.

For detailed instructions, please refer to the documentation.

## API Documentation

The API documentation is available using Swagger/OpenAPI. You can access it locally at `http://localhost:3000/swagger/index.html` when the API Gateway is running.

## License

This project is open-source and available under the [MIT License](LICENSE).
