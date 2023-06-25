#!/bin/sh
while true; do
    echo "Opening server in 127.0.0.1:5000"
    nc -t -s 127.0.0.1 -l -p 5000
    sleep 1
    echo "Connection terminated..."
done