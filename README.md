# Simple Test Data Generator

Simple app generating test data based on JSON configuration file. Uses [fake]
golang module.

## Getting

    go get github.com/pavel-popov/stdg

## Usage

    stdg -schema config.json -rows 10000 -columns a,b,c


## References

* [fake]


[fake]: https://github.com/icrowley/fake
