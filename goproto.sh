#!/usr/bin/env bash

cd proto;
protoc --go_out=plugins=grpc:../lib/proto-go *.proto

