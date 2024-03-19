## Braz E-commerce
This project is a small e-commerce done with good practices and the best of golang!

### First steps
As a new member of the engineering team the first step that you'll need to do is to execute the project setup command:
```sh
make setup
```

This project is developer friendly and will watch your changes and automatically reload the microservices. In order to achieve this result, you'll need to execute:
```sh
make start
```
The command above will start all microservices and infrastructure resources in docker containers and start to scan for file changes. To stop the running containers you can easily execute:
```sh
make stop
```

### Unit tests
#### Running unit tests
There's also a command to execute unit tests and can be easily invoked through make command:
```sh
# running unit tests:
make test

# running unit tests with coverage report:
make coverage
```

#### Creating interface mocks
The first step to start creating interface mocks is to install mockgen package
```sh
go install go.uber.org/mock/mockgen@latest
```
The next step is to go to the file where the target interface is located and copying the relative path and running the following command:
```sh
mockgen -source=${INSTANCE_FILE_RELATIVE_PATH} -destination=services/${SERVICE_NAME}/mocks/${INTERFACE_NAME}_interface_mock.go -package=mocks -write_generate_directive
```
