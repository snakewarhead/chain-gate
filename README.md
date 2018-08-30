# chain-gate

## summy

This is a lightweight gateway for all blockchain programs, which provides a common interface to mask all blockchain program differences and return consistent data formats.

## install

1. install sqlite3

    ```sh
    sudo apt-get install pkg-config
    sudo apt-get install sqlite3
    sudo apt-get install libsqlite3-dev

    # sudo apt-get install sqlitebrowser
    ```

1. enable coin in the table 'coin' of database 'chain_data.db'

1. ~~install libs~~

    ```sh
    go get github.com/jeanphorn/log4go
    go get github.com/mattn/go-sqlite3
    go get github.com/shopspring/decimal
    ```

1. install

    ```sh
    go get github.com/snakewarhead/chain-gate
    ```