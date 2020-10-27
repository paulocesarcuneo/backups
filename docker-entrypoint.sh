#!/bin/bash

if [[ -z "$SERVER" ]]; then
    ./server
else
    ./client
fi
