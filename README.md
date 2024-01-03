# Key Value DB (Redis) in Go

The is a command-line utility that allows you to interact with a simple in-memory key-value database through a TCP server.

## Installation

You must have Go(1.18) installed on your machine. Follow the steps below to install and set up the tool:

1. Clone the repository.
2. Run the following command in the project directory to build the CLI tool:

   ```shell
   go build -o db
   ```

   This will create an executable file named `db` in the project directory.

## Usage

To start the TCP server, follow these steps:


1. Create a `.env` file in the root directory and set the required values to the below environment variables. We have provided a `.env.example` file in the root directory for reference.

    1. Ensure the environment variable `TCP_PORT` is set to the desired port number on which the TCP server should listen. For example, you can set it to `8003` by running:

       ```shell
       export APP_PORT=8003
       ```

    2. Set the environment variable to set the number of in-memory databases (`DB_COUNT`). For example:

       ```shell
       export DB_COUNT=16
       ```
       
2. Run the following command to start the TCP server:

   ```shell
   ./db
   ```

3. The TCP server will start and display a message indicating it running.

4. Open another terminal and use a tool like `nc` or `telnet` to connect to the TCP server. For example:

   ```shell
   nc localhost 8003
   ```
   
5. Once connected, you can interact with the CLI tool by entering commands. A `>` symbol denotes the command prompt.

6. The available commands are case-insensitive and can be entered in the following format:

   ```
   COMMAND [argument1] [argument2] ...
   ```

   Replace `COMMAND` with a supported command and provide the necessary arguments.

7. The CLI tool supports the following commands:

    - `SET key value`: Sets the value of the specified key in the current database.
    - `GET key`: Retrieves the value of the specified key from the current database.
    - `DEL key`: Deletes the specified key from the current database.
    - `INCR key`: Increments the value of the specified key by 1.
    - `INCRBY key increment`: Increments the value of the specified key by the specified increment.
    - `MULTI`: Starts a transaction block.
    - `EXEC`: Executes all commands in a transaction block.
    - `DISCARD`: Discards all commands in a transaction block.
    - `COMPACT`: Compacts the database by removing expired keys.
    - `SELECT` index: Switches to the specified database index (0-based).
    - `DISCONNECT` disconnect the connected client from the TCP server.

   Replace key, value, index, and increment with the appropriate values.

8. After entering a command, the CLI tool will display the command result.

9. To exit the CLI tool, close the `nc` connection or terminate the terminal session or use the `DISCONNECT` command.

## Known issues

1. The `ui.getDbIndex(result any)` function may return a wrong database index (i.e zero) `result any = []domain.DBResult` and that slice is empty.

## TODO

1. Fix ISSUE #1
2. Add fast persistent storage, not just in-memory.

## License
This project is licensed under the [MIT License](./LICENSE)