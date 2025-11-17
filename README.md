# Mini-Processor Takehome Project

This README will be filled in later; removed the initial how-to instructions for brevity and initial commit.

## Initial Setup

Please run:

1. `go mod download` to fetch dependencies
2. To use the `ginkgo` command line tool to run tests, please run `go install github.com/onsi/ginkgo/v2/ginkgo`

## Testing

The project script provides the ability to run all testing (including integration tests) via: `./project test`

This will start the pubsub emulator, topic creation and subscription creation docker images.
Once these have started, the unit tests and integration tests will run together.

The coverage data file will be generated in the `<repository_root>/.reports/coverage.out` file.

## HTML Coverage Report

To check the coverage for the code via the standard golang HTML coverage report (again including integration tests), run: `./project coverage`

This will start the pubsub emulator, topic creation and subscription creation docker images.
Once these have started, the unit tests and integration tests will run together and present the HTML coverage on completion.

