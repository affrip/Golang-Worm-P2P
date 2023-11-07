#!/usr/bin/bash
CGO_ENABLED=0 go build -o bot -a -ldflags '-extldflags "-static"' . 
