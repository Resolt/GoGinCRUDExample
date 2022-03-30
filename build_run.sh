#!/bin/sh

export DB_HOST="localhost"
export DB_NAME="ggce"
export DB_USER="ggce"
export DB_PASS="ggce"
export DB_PORT="5432"
export PORT="8000"
export AMQP_USER="guest"
export AMQP_PASS="guest"
export AMQP_HOST="localhost"
export AMQP_PORT="5672"
export AMQP_VHOST=""
export AMQP_EXCHANGE="ggce"
export AMQP_QUEUE="ggce"

go build -o ggce && ./ggce