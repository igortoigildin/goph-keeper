# Password Manager GophKeeper

## Project Description

The project implements a storage system for storing different types of data, such as passwords, text, and binary data.

## Features

1. **Storage System**

   - Before using this CLI app, every user should register and log in with their credentials.
   - Using GophKeeper, users can save their passwords, text, bank card details, or any type of binary files.
   - After uploading their data, each user will receive a unique ID to download their data later.

2. **Technologies**
   - **Go** — backend language for the application’s business logic.
   - **PostgreSQL** — database for storing workshop and booking information.
   - **pgx** — library used for database operations.
   - **grpc** — framework used for communication between server and cli client.
   - **Cobra** — library used for building cli commands.
   - **Viper** — package used for managing configuration.
   - **Minio** — object storage system that is compatible with the Amazon S3 (Simple Storage Service) API.

## Repository Structure

- `cmd/server/main.go` — application server entry point.
- `cmd/client/main.go` — application client entry point.
- `internal/server` — core application logic for server.
  - `app` — DI container for server.
  - `storage` — module for database interactions and schema handling.
  - `api` — app server handlers.
  - `models` — ubiquitous domain server models.
- `internal/client` — core application logic for client.
  - `/grpc/app` — cobra commands initialization for client.
  - `/grpc/service` — client service functions.

## Getting Started

### Requirements

- **Go**
- **PostgreSQL**
- **Docker**

### Installation and Run

1. Clone the repository:

   ```bash
   git clone https://github.com/igortoigildin/goph-keeper.git
   cd goph-keeper
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```
3. Start docker-compose:

   ```bash
   docker-compose up -d
   ```

4. Start the server:
   ```bash
   go run cmd/server/main.go
   ```
5. Open new terminal and build the client:

   ```bash
   make build-client
   ```

6. Install and get all necessary dependencies:

   ```bash
   make install-deps
   make get-deps
   ```

7. Generate all api for gRPC:

   ```bash
   make generate
   ```

8. Roll migrations:

   ```bash
   make migration-up
   ```

### Commands Examples

#### Registration and Login

1. Registration. Open separate terminal and run (docker and server must be running):

```bash
    bin/client create user -l temp_login -p 123
```

2. Login.

```bash
    bin/client login user -l temp_login -p 123
```

#### Save commands

```bash
    bin/client save text -t sample_text
    bin/client save password -l fakelogin -p 1234
    bin/client save bin -n migration -p migration.sh
    bin/client save card -n 1234 -c 329 -e 12/04/1005
```

#### Download commands

Your should use your unique id for data to make downloads as per below examples.

```bash
    bin/client download password -i cb4b3e82-fd37-4faa-8d38-603a65990a57
    bin/client download text -i 80a0766f-1b53-4603-9738-1bcc954fa4d3
    bin/client download card -i 0d9efc10-36cc-425b-a599-465b21855977
    bin/client download bin -i 092049f9-2719-44eb-aa12-25e167dcba13 -n samplename
```
