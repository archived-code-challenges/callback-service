# Golang Callback Service

## Requirements

- Golang
- Docker
- Docker compose

## Instructions to run the project

If the host operating system is MacOS, I've tried to configure the whole thing with one command ([`make compose-mac`]((Makefile#L46))) but couldn't make work a docker-compose approach [due to networking limitations](https://docs.docker.com/docker-for-mac/networking/#known-limitations-use-cases-and-workarounds) when talking between containers.

So basically, to **start everything on a MacOS** you must have golang installed:

    make compose-local-mac

> Any flag-attribute listed when the application starts can be changed at execution time, check [Makefile (line 46)](Makefile#L46) for an extended ilustration.
> Therefore, the client server could be changed if required:
> `make compose-local-mac LDFLAGS=--callback-service-address=http://1234host:7777`

To start the client service that will query the server each 5 seconds and also provide the status of the objects:

    make client-start

Probably with a Linux machine, the docker-compose approach will work, setting up everything at once:
*(Not tested)*

    make compose-linux

To start the client service that will query the server each 5 seconds and also provide the status of the objects:

    make client-start

After a successful execution, the service should be running on port 9090. Hit the endpoint to check that all is running as expected:

- [:9090](http://localhost:9090/)

Tracing can be found at:

- [:9411/zipkin/](http://localhost:9411/zipkin/)
- [:6060/debug/pprof/](http://localhost:6060/debug/pprof/)

To explore what else you can do with the Makefile:

    make

## Testing

Tests point to a different database that will be created just for test purposes.

After that, the project tests can be run by typing:

    make test

Integration tests (which require both services running) can be run by typing:

    make test-integration

## Packaging

    .
    ├── cmd                     # Entrypoint
    │   ├── callback-service    # Main API
    │   └── client-service      # client callbacks API
    ├── doc                     # Documentation, images and helpful files
    └── internal
        ├── handlers            # HTTP layer & integration tests
        ├── middleware
        ├── models              # Business logic
        └── web                 # Framework for common HTTP related tasks
