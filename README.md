# GotCC - A Transaction Consistency Framework

GotCC is a Go-based framework designed to manage distributed transactions with a focus on achieving eventual consistency. It utilizes a layered architecture to separate concerns, making it easier to maintain and extend.

## Project Structure

- **cmd/gotcc/main.go**: Entry point of the application. Initializes the application, sets up configurations, and starts the engine.
- **internal/config**: Contains configuration loading and schema definition.
  - **loader.go**: Functions to load and parse YAML/JSON configuration files into memory.
  - **schema.go**: Defines the structure of the configuration schema used for validation and loading.
- **internal/engine**: Implements the core logic of the transaction engine.
  - **engine.go**: Manages the overall flow of the transaction process.
  - **state_machine.go**: Defines the state machine for managing different states of the transaction process.
  - **scheduler.go**: Handles the scheduling of actions and tasks within the transaction flow.
  - **retry.go**: Implements the retry logic for tasks that fail during execution.
- **internal/executor**: Defines action executors for executing tasks.
  - **executor.go**: Interface for action executors and manages task execution.
  - **plugins**: Contains implementations of various action executors.
    - **http_executor.go**: HTTP action executor for making HTTP requests.
    - **db_executor.go**: Database action executor for executing SQL commands.
    - **mq_executor.go**: Message queue action executor for sending messages to a queue.
- **internal/persistence**: Manages database interactions.
  - **store.go**: Manages the database connection and transaction handling.
  - **repository.go**: Defines the repository pattern for accessing and manipulating task and instance data.
  - **migrations/def.sql**: SQL definitions for creating necessary database tables.
- **internal/dao**: Provides data access methods.
  - **task_dao.go**: Data access methods for task-related operations.
  - **instance_dao.go**: Data access methods for transaction instance-related operations.
- **internal/model**: Defines data structures for managing flows, instances, and tasks.
  - **flow.go**: Data structures and methods for managing flow definitions.
  - **instance.go**: Data structures and methods for managing transaction instances.
  - **task.go**: Data structures and methods for managing tasks.
- **pkg/api/http**: Sets up the HTTP server and routes for the API.
  - **server.go**: Sets up the HTTP server.
  - **handlers.go**: Defines HTTP handlers for processing requests.
- **pkg/transport**: Implements client-side logic for interacting with external services.
  - **client.go**: Client-side logic for service interaction.
- **configs**: Contains example configuration files.
  - **example.yaml**: Example configuration in YAML format.
  - **example.json**: Example configuration in JSON format.
- **scripts**: Contains utility scripts.
  - **migrate.sh**: Script for running database migrations.
- **deployments**: Contains deployment configurations.
  - **docker-compose.yml**: Docker Compose configuration for deployment.
- **test**: Contains tests for the application.
  - **integration**: Integration tests.
  - **unit**: Unit tests.
- **go.mod**: Module dependencies and versioning.
- **Makefile**: Build and deployment commands.

## Getting Started

1. **Clone the repository**:
   ```
   git clone <repository-url>
   cd gotcc
   ```

2. **Install dependencies**:
   ```
   go mod tidy
   ```

3. **Run the application**:
   ```
   go run cmd/gotcc/main.go
   ```

4. **Configuration**: Modify the configuration files in the `configs` directory to suit your needs.

## Usage

- The framework supports defining transaction flows and executtno略、补偿回滚与幂等保障。
rdedns with retry and compensation mechanisms.
- You can define your tno略、补偿回滚与幂等保障。
rdedn the configuration files and implement custom action executors as no略、补偿回滚与幂等保障。
rded.

## Contributing

Contributions are welcome! Please open an issue o略、补偿回滚与幂等保障。
r submit a pull request for any enhancements or bug fixes.

## License

This 略、补偿回滚与幂等保障。
project is licensed under the MIT License. See the LICENSE file for more details.略、补偿回滚与幂等保障。
