#!/bin/sh
# Development test server using netcat.
# Uses ncat (nmap) if available, otherwise plain nc.
# Note: some nc versions (e.g. OpenBSD) use 'nc -l 127.0.0.1 5000'

if command -v ncat >/dev/null 2>&1; then
    NC=ncat
    LISTEN="127.0.0.1 5000"
else
    NC=nc
    LISTEN="-l -s 127.0.0.1 -p 5000"
fi

while true; do
    echo "Opening server on 127.0.0.1:5000"
    $NC $LISTEN
    sleep 1
    echo "Connection terminated..."
done