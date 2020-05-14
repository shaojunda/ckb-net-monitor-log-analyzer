# CKB Net Monitor Log Analyzer
[![License](https://img.shields.io/badge/license-MIT-green)](https://github.com/shaojunda/ckb-net-monitor-log-analyzer/blob/master/LICENSE)
[![Go version](https://img.shields.io/github/go-mod/go-version/shaojunda/ckb-net-monitor-log-analyzer)](https://github.com/moovweb/gvm)

CKB Net Monitor Log Analyzer is a log analyzer for analyzing [CKB Net Monitor](https://github.com/quake/ckb-net-monitor) Logs. So far we can analyze **Block Propagation Delays** and **Transaction Propagation Delays** from the log.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. 

### Prerequisites

* [Golang](https://golang.org) 1.14.2 and above

* [PostgreSQL](https://www.postgresql.org/) 9.4 and above

### Installing

`go get -v github.com/shaojunda/ckb-net-monitor-log-analyzer`

### Basic Usage

1. Edit `config.yaml` file to set various parameters

    ```yaml
    monitor_log_file_path: run1.log  # The location of the file to be analyze
    process_count: 10000 # Save data to DB after processing {process_count} data
    pg_host: localhost
    pg_port: 5432
    pg_user: postgres
    pg_password: postgres
    pg_db_name: db_name
    ```

2. Create two tables named `block_propagation_delays` and `transaction_propagation_delays`

3. Start analysis
`go run main.go`

During the analysis, the program can be interrupted the next time it will be executed at the last execution position.


## Contributing
CKB Net Monitor Log Analyzer is an open-source project, and thank you very much for your contribution. Please check out [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines about how to proceed.


## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
