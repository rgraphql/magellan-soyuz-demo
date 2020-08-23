# Magellan/Soyuz Demo

> Demonstration of Go and TypeScript two-way graphql over WebSockets.

## Introduction

This is a demonstration of [Magellan] and [Soyuz] communicating over a
WebSocket. It uses React components to demonstrate the features of the system.

[Magellan]: https://github.com/rgraphql/magellan
[Soyuz]: https://github.com/rgraphql/soyuz

## Getting Started

In the project directory, you can run:

### `magellan analyze`

```bash
cd ./server
./compile.bash
```

This will analyze the code and re-generate the resolvers.

### `yarn start` and `server`

First start the Go server:

```bash
cd ./server
go build -v
./server
```

Next run: `yarn start`

Runs the app in the development mode.<br />
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.

You will also see any lint errors in the console.

The demo is a JSON view of the current resolver output of the following query:

```graphql
{
  counter
  names
  allPeople {
    name
    height
  }
}
```

If you just see {} make sure you have run the demo server. You may need to
refresh a few times. The demo code is quite brittle at the moment.

### `yarn test`

Launches the test runner in the interactive watch mode.

### `yarn build`

Builds the app for production to the `build` folder.

