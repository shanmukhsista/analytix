#!/bin/bash
protoc -I./ --go_out=../generated --go_opt=paths=source_relative ./models.proto
