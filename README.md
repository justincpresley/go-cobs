<div align="center">

![Visual](/docs/README_VISUAL.png)

[![Test](https://img.shields.io/github/actions/workflow/status/justincpresley/go-cobs/test.yaml?branch=master&label=Test)][1]
[![CodeQL](https://img.shields.io/github/actions/workflow/status/justincpresley/go-cobs/codeql.yml?branch=master&label=CodeQL)][2]
[![CodeFactor](https://img.shields.io/codefactor/grade/github/justincpresley/go-cobs/master?label=CodeFactor)][3]
[![Language](https://img.shields.io/github/go-mod/go-version/justincpresley/go-cobs/master?label=Go)][4]
[![Version](https://img.shields.io/github/v/tag/justincpresley/go-cobs?label=Latest%20version)][5]
[![Commits](https://img.shields.io/github/commits-since/justincpresley/go-cobs/latest/master?label=Unreleased%20commits)][6]
[![License](https://img.shields.io/github/license/justincpresley/go-cobs?label=License)][7]

</div>

***go-cobs*** is a [Go](https://go.dev/) library implementing *Consistent Overhead Byte Stuffing* (COBS) functionality.

*What is Consistent Overhead Byte Stuffing?* [Wiki](https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing) - [Technical Paper](http://www.stuartcheshire.org/papers/cobsforton.pdf)

The goal of 'COBS' is to remove a special byte within given data by replacing all special bytes with "flags", a byte telling where the next special byte is. There is minimal overhead in COBS as indicated by the paper itself.

## Usage

It is suggested to try out this library by running ```go run .``` instead of one the [examples](https://github.com/justincpresley/go-cobs/tree/master/examples) to visualize **go-cobs**.

The full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/justincpresley/go-cobs).

Additionally, [usage](https://github.com/justincpresley/go-cobs/blob/master/docs/USAGE.md) outlines everything you need to know about types, config, and anything that may be unclear outside of the API.

## License

***go-cobs*** is an open source project licensed under ISC. See LICENSE.md for more information.

[1]: https://github.com/justincpresley/go-cobs/actions/workflows/test.yaml
[2]: https://github.com/justincpresley/go-cobs/actions/workflows/codeql.yml
[3]: https://www.codefactor.io/repository/github/justincpresley/go-cobs
[4]: https://go.dev/
[5]: https://github.com/justincpresley/go-cobs/releases
[6]: https://github.com/justincpresley/go-cobs/compare/v1.5.0-alpha.1...HEAD
[7]: https://en.wikipedia.org/wiki/ISC_license
