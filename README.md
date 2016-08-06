# Simple Test Data Generator

Simple app generating test data based on JSON configuration file. Uses [fake]
golang module.

## Building

    go build

## Usage

    stdg -config config.json -rows 10000

[fake]: https://github.com/icrowley/fake
