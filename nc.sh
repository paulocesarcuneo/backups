#!/bin/bash


docker run --rm -i --network=$1 backups nc ${@:2}
