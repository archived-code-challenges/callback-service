# Golang Callback Service

<img src="https://img.shields.io/badge/status-development-yellowgreen"/></a>

Golang Callback Service is a simple API that listens for `POST` requests in `/callback` and requests object information for each incoming `object_id` concurrently. Validates the request, returns a response, and processes the payload. Thus, avoiding bottlenecks in a thread-safe way.

## Dependencies

This project was designed against Go 1.16.
Supporting services like the database are hosted in Docker and use Docker Compose to startup.

- Golang
- Docker
- Docker compose

## API Request

`cmd/callback-service`

| Endpoint        | HTTP Method   | Description         |
| --------------- | :-----------: | :-----------------: |
| `/callback`     | `POST`        | `Create objects`    |
| `/`             | `GET`         | `Health check`      |

`cmd/client-service`

| Endpoint        | HTTP Method   | Description         |
| --------------- | :-----------: | :-----------------: |
| `/objects/:id`  | `GET`         | `Retrieves object`  |

## Decisions taken on this implementation

### Product

- Given a request to the server, the payload will be validated and a response will be returned to the client before processing is completed.
- Errors that occur while data is being processed in goroutines are logged.
- Deletion of data on the database is executed exactly N seconds after insertion.
- Unit and integration tested.

### Design

- Unhandled error fail-safe. An unhandled error will output as internal_error.
- Captures logs. Stores traces with Zipkin and profiles memory usage with *pprof*.
- The core parts of the application are documented.
- Service variables can be modified at execution time.
- The folder structure follows a package-oriented design to maintain the best possible purpose, usability, and portability.
- The project dependencies are located in the vendor folder.
- When stopped, the service tries to shut down gracefully. Logging any errors that may occur.

## Instructions to run the project

If the host operating system is MacOS, I've tried to configure the whole thing with one command ([`make compose-mac`]((Makefile#L46))) but couldn't make work a docker-compose approach [due to networking limitations](https://docs.docker.com/docker-for-mac/networking/#known-limitations-use-cases-and-workarounds) when talking between containers.

Thus, to **start everything on a MacOS** you must have the previously listed [dependencies](#dependencies) installed:

    make compose-local-mac

> Any flag attribute listed when the application starts can be changed at execution time. Check [Makefile (line 46)](Makefile#L46) for an extended illustration.
> Therefore, the client-server could be changed if required:
> `make compose-local-mac LDFLAGS=--callback-service-address=http://1234host:7777`
> This functionality can be used to easily (for example) switch URL references between the services involved.

With a Linux machine, the docker-compose approach should work, setting up everything at once:
*(Not tested)*

    make compose-linux

To start the client service that will query the server every 5 seconds and also provide the status of the objects:

    make client-start

To explore what else you can do with the Makefile:

    make

After successful execution, the service should be running on port 9090. Hit the *health* endpoint to check that all is running as expected:

- [:9090](http://localhost:9090/)

Tracing can be found at:

- [:9411/zipkin/](http://localhost:9411/zipkin/)
- [:6060/debug/pprof/](http://localhost:6060/debug/pprof/)

![zipkin-traces](/doc/zipkin-traces.png)

## Testing

The tests point to a different database created just for them, so they are isolated from the core service.

Project tests can be run by typing:

    make test

Integration tests (which require both services running) can be performed by typing:

    make test-integration

## Packaging

    .
    ├── cmd                     # Entrypoint
    │   ├── callback-service    # Main API
    │   └── client-service      # Client API (provided by the task)
    ├── doc                     # API Documentation, images, and helpful files
    └── internal
        ├── handlers            # HTTP layer & integration tests
        ├── middleware
        ├── models              # Business logic
        └── web                 # Framework for common HTTP related tasks

## Use of external libraries

- [gorm](https://github.com/go-gorm/gorm). ORM library.
- Zipkin, go-zipkin, and Opencensus. Distributed tracing system
- [go-chi](https://github.com/go-chi/chi). Router for building Golang HTTP services.
- [conf](github.com/ardanlabs/conf). Support for using environmental variables and command-line arguments for configuration.

## Author

Noel Ruault - [@noelruault](https://github.com/noelruault)

## License

Copyright © 2020 [@noelruault](https://github.com/noelruault).
This project is [MIT](https://opensource.org/licenses/MIT) licensed.
