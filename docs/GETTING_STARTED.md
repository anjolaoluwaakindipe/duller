## GETTING STARTED

- The main program of this project can be found in `./cmd/duller/main.go`. In order to run the code as a CLI tool you can use the following command

```bash
go run ./cmd/duller/main.go [COMMAND] [ARGS]
```

- There are a currently two sub commands involved with duller `disc` and `gate`. `disc` is for the service discovery server and `gate` is for starting up the gateway.

- To get a list of flags that can be set for both commands you can pass the `-h` or `--help` flag.

```bash
go run ./cmd/duller/main.go disc -h
```
