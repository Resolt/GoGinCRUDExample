#!/bin/sh

export DB_HOST="localhost"
export DB_NAME="ggce"
export DB_USER="ggce"
export DB_PASS="ggce"
export DB_PORT="5432"
export PORT="8000"

go build -o ggce && ./ggce