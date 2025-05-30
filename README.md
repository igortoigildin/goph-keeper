# Password Manager GophKeeper

## Project Description

The project implements a storage system for storing different types of data, such as passwords, text, and binary data.

## Features

1. **Storage System**

   - Before using this CLI app, every user should register and log in with their credentials.
   - Using GophKeeper, users are able to save their passwords, text, bank card details, or any type of binary files.
   - After uploading their data, each user will receive a unique ID to download their data later.

2. **Technologies**
   - **Go** — backend language for the application’s business logic.
   - **PostgreSQL** — database for storing user's information.
   - **pgx** — library used for database operations.
   - **grpc** — protocol used for communication between server and cli client.
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
  - `/grpc/models` — ubiquitous domain client models.
  - `/grpc/storage` — module for client local database interactions.

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

7. Generate api:

   ```bash
   make generate
   ```

8. Roll migrations:

   ```bash
   make migration-up
   ```

9. Run tests (docker and server must be running):
   ```bash
   make test-cover
   ```
10. Generate TLS certs:

```bash
make certs
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

#### Save and download text data

Please note, that you should use your unique id for data to make downloads.
You also my add your optional additional metadata for your secrets with flags -i or -d. Please see examplse below.

1. Save text data

```bash
    bin/client save text -t sample_text -i optinal_metadata
```

2. Download text data.

```bash
   bin/client download password -i cb4b3e82-fd37-4faa-8d38-603a65990a57
```

3. Save login&password data

```bash
    bin/client save password -l fakelogin -p 1234 -d optinal_metadata
```

4. Download login&password data.

```bash
   bin/client download password -i cb4b3e82-fd37-4faa-8d38-603a65990a57
```

5. Save bank details

```bash
    bin/client save card -n 1234 -c 329 -e 12/04/1005  -i optinal_metadata
```

6. Download bank details.

```bash
    bin/client download card -i 0d9efc10-36cc-425b-a599-465b21855977
```

7. Save binary data

```bash
    bin/client save bin -n migration -p migration.sh -i optinal_metadata
```

8. Download binary data. The file will be downloaded to 'client_files' directory.

```bash
    bin/client download bin -n tempname -i 092049f9-2719-44eb-aa12-25e167dcba13
```

9. List all secrets saved

```bash
    bin/client list all
```

10. Sync secrets data between client and server

```bash
    bin/client sync all
```
