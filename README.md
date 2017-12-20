# Go Flipper

Go Flipper is a port of [Flipper](https://github.com/jnunemaker/flipper) to Go.

It's designed to be compatible at the storage level with the Ruby gem, so you can toggle features using any of them and access them with the other.

This project is in an early stage and doesn't include all the features that the Ruby gem includes at the moment.

Caveats:

- This implementation assumes that all actors have a FlipperID method that returns a string. It doesn't work with other id types at the moment.

## Installation

```
go get github.com/calavera/go-flipper
```

## Usage

See the Godoc for examples and API information:

https://godoc.org/github.com/calavera/go-flipper

## License

[MIT](LICENSE)
