# Merkle Challenge

Solution to the take-home assignment of an interview process I took part in 2022. Description of the assignment can be found [here](problem/README.md).

## Test plan
The test plan can be found [here](docs/Test_Plan.md).

## Development choices
Despite having worked in the past with some of the languages in use at your company (Java, Python, and C#), I chose to write the tests in Golang (Go) because it is the language I use on a daily basis since 2020. For this challenge, Go has been used together with a testing framework called [Ginkgo](https://onsi.github.io/ginkgo/) and its matcher library [Gomega](https://onsi.github.io/gomega/).

The tests are written in a BDD-like flavor: each test case description is wrapped into its own [It](https://onsi.github.io/ginkgo/#spec-subjects-it) function. I think this approach eases reading the test. To further ease the process, I intentionally left some duplicated code among the tests (so the reader can identify quickly what each test does).

### Go code folder structure
[suites](./suites) contains the entry point of the program and the test cases.

[internal](./internal) contains the components the tests delegate tasks to (i.e. encoding, HTTP client, hash calculation etc.)

## How to run the tests
The test cases 1, 2, 3, 9, 10, 11 and 12 of the [test plan](docs/Test_Plan.md#test-cases-document) have been automated.

You have two options for running the tests: Docker **or** Golang and Java. For the latter, you also need the Ginkgo binary. You can install it with the command:
```bash
go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo
```

### Docker
The repo contains two Dockerfiles and a docker-compose:
- [server](./docker/server/Dockerfile) copies the `.jar` along with two files in a Debian distribution with openjdk 18 pre-installed
- [tests](./docker/tests/Dockerfile) builds the tests binary and copies it together with the [config file](./docker/tests/merkle.yml) in a light distribution (Linux Alpine)
- [compose](./docker/docker-compose.yml) builds the two custom images and spawns the containers. The tests container is started only when the server container is healthy (i.e. the .`jar` is running)

To run the tests, run the following command in the terminal:
```bash
make dev
```
The docker-compose will be launched in foreground (you must press "ctrl+c" or "cmd+c" to regain control of the window). In this instance, I prefer this approach to running the command in background and then having to run `docker-compose logs -f` to see the execution status.

There are two commands more:
```bash
make nodev # stop the containers
make devrebuild # rebuild the images
```

### Go and Java
Before running the tests you might want to execute `go mod tidy` to download potentially missing dependencies.

`make server`: spawns the web server (`.jar`)

`make tests`:  runs the tests with the Ginkgo binary and uses [this](./merkle.yml) config file. You can also provide a custom config file to the make target (`config.yml` is ignored by git). Example:
```bash
CFG=../config.yml make tests # ".." because the tests are in the "suites" folder
```

The config file is also used to make some of the tests data-driven.

## Tools
- [trivy](https://github.com/aquasecurity/trivy) has been used to identify vulnerabilities in the Dockerfiles
- [hadolint](https://github.com/hadolint/hadolint) has been used to lint the Dockerfiles
- [golangci-lint](https://golangci-lint.run/) has been used to lint the Go code

## System spec
All the work has been done on a 2021 16" Apple Macbook Pro with the following specs:
- M1 Pro processor (10-core CPU and 16-core GPU)
- 32 GB of RAM
- 1 TB SSD
### Software
- macOS Monterey 12.4
- Go 1.18.3
- Ginkgo 2.1.4
- JDK 18.0.1.1

## Web server errors
1. When supplying a large file (6 GB or more) to the `.jar`, the process exits immediately with the message: `java.lang.OutOfMemoryError: Required array size too large`. Running the jar with increased heap size (tested with `-Xms10G -Xmx16G`) did not solve the issue.

2. When supplying a < 1 KB file to the web server, a call to the endpoint `/piece/:hash/0` returns: `{404, "No such route"}`. The app log prints: `java.lang.ArrayIndexOutOfBoundsException: Index 1 out of bounds for length 1`. You can reproduce the error by running the following commands:
```bash
java -jar problem/merkle-tree-java.jar problem/small
# The command below in a separate terminal window
curl localhost:8080/piece/c4dbc298900101df82899b1b9dd2c9752f3dc59c8d79843c9083210cf882040e/0
```
