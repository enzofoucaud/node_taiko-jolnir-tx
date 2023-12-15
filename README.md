# taiko-jolnir-tx

## Description

This is a simple tool to make it easier to do transactions for Node on the [Taiko](https://taiko.xyz/) platform.

## Installation

Install [Golang](https://go.dev/doc/install). Then run:

```bash
git clone github.com/enzofoucaud/node_taiko-jolnir-tx
```

## Usage

You can use built version of the tool or build it yourself.

Copy the `config.example.json` file to `config.json` and fill in the values.

```bash
cp config.example.json config.json
```

### Build

```bash
go build . -o taiko-jolnir-tx
```

### Run

```bash
./taiko-jolnir-tx
```

Example:

```bash
./taiko-jolnir-tx
```
