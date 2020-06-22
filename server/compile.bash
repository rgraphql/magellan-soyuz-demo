#!/bin/bash

go run \
   github.com/rgraphql/magellan/cmd/magellan \
   analyze --schema ../schema.graphql \
   --go-pkg github.com/rgraphql/magellan-soyuz-demo/server/simple \
   --go-query-type RootResolver \
   --go-output ./resolve/resolve_generated.go
