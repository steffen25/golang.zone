#!/usr/bin/env bash

go test -v -race ./... -coverprofile=cover.out && go tool cover -html=cover.out