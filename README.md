# trade-srv

> A client for distributed financial news webscraping.

### Usage

1. Download the [latest release](https://github.com/sshh12/trade-srv/releases)
2. Have a postgres server
3. Register symbols using `$ ./trade-srv-index -pg_url postgres://username:password@127.0.0.1:5432/tradesrv -add_sym TSLA,AAPL,AMZN`
4. Start scraping with `$ ./trade-srv-index -pg_url postgres://username:password@127.0.0.1:5432/tradesrv -run_all_events -log info`, use `$ ./trade-srv-index -h` for all options.
