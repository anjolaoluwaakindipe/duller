# DULLER

This is a service discovery and gateway application for microservices inspired by Eureka and Zuul. It makes use of tcp calls to registers microservices and an api gateway to proxy users to appropriate services by making special calls to the discovery server.g

## GETTING STARTED

1.  This project uses make files. In order to run the code as a CLI tool you can use the following command

```bash
make disc
```

2. A list of flags can be gotten with

## FURTHER PLANS

1. Create multiple clients for other programming languages. Will need help contributing
2. Add health monitoring and potential logging service statuses
