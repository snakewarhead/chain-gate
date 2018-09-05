# chain-gate

## summary

This is a lightweight gateway for all blockchain programs, which provides a common interface to mask all blockchain program differences and return consistent data formats.

## install

1. install go
    > go version go1.10.3 linux/amd64 or above

1. install gate

    ```sh
    go get github.com/snakewarhead/chain-gate
    ```

1. install sqlite3

    ```sh
    sudo apt-get install pkg-config
    sudo apt-get install sqlite3
    sudo apt-get install libsqlite3-dev

    # sudo apt-get install sqlitebrowser
    ```

1. enable and config 'coin' in the table 'coin' of database 'chain_data.db'
