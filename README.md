# trade-srv

> A client for distributed financial news webscraping.

![img](https://user-images.githubusercontent.com/6625384/102290101-ebe56d80-3f05-11eb-88e1-1801c3383cd5.png)

## Features

- [Several](https://github.com/sshh12/trade-srv/tree/main/indexers) financial news sites
- Plug-and-Play webscraping
- Request optimizations
- Increase timestamp accuracy by running more clients

## Usage

1. Download the [latest release](https://github.com/sshh12/trade-srv/releases).
2. Have a fresh PostgreSQL database.
3. Register symbols using `$ ./trade-srv-index -pg_url postgres://username:password@127.0.0.1:5432/tradesrv -add_sym TSLA,AAPL,AMZN`
4. Run:

```bash
$ ./trade-srv-index -h
$ ./trade-srv-index -pg_url postgres://username:password@127.0.0.1:5432/tradesrv -run_all_events -log info
```
