#!/bin/bash
# It is a simple bash script to launch all components. It could be better
# implemented using Python and have a more thorough control of the processes,
# but it is intended for rapid demonstration only.
# CMake could be used for building, but again, the intention is just to get
# something working without adding overhead of installing new tools.

# config

export IP_VERSION=4

PAYMENT_API_SERVICE_DIR=api
export PAYMENT_API_SERVICE_HOST=0.0.0.0
export PAYMENT_API_SERVICE_PORT=8080

MERCHANT_SERVICE_DIR=merchant
export MERCHANT_SERVICE_HOST=0.0.0.0
export MERCHANT_SERVICE_PORT=50051

RATE_LIMITER_SERVICE_DIR=ratelimiter
export RATE_LIMITER_SERVICE_HOST=0.0.0.0
export RATE_LIMITER_SERVICE_PORT=50052

LEDGER_SERVICE_DIR=ledger
export LEDGER_SERVICE_HOST=0.0.0.0
export LEDGER_SERVICE_PORT=50053

BANK_SIMULATOR_DIR=bank
export BANK_SIMULATOR_HOST=0.0.0.0
export BANK_SIMULATOR_PORT=8000

MERCHANT_UI_DIR=merchant-ui
export MERCHANT_UI_HOST=0.0.0.0
export MERCHANT_UI_PORT=3000
export REACT_APP_PAYMENT_GATEWAY_HOST=0.0.0.0
export REACT_APP_PAYMENT_GATEWAY_PORT=8080

# building components
cd $MERCHANT_SERVICE_DIR && \
    go build -o app && \
    cd ..
cd $RATE_LIMITER_SERVICE_DIR && \
    go build -o app && \
    cd ..
cd $LEDGER_SERVICE_DIR && \
    go build -o app && \
    cd ..
cd $BANK_SIMULATOR_DIR && \
    rm -rf venv && \
    python3 -m venv venv && \
    source venv/bin/activate && \
    pip install --upgrade pip && \
    pip install -r requirements.txt && \
    cd ..
cd $PAYMENT_API_SERVICE_DIR && \
    go build -o app \
    && cd ..
cd $MERCHANT_UI_DIR && \
    rm -rf node_modules package-lock.json && \
    npm install &&
    cd ..

# start components
cd $MERCHANT_SERVICE_DIR && ./app &
cd $RATE_LIMITER_SERVICE_DIR && ./app &
cd $LEDGER_SERVICE_DIR && ./app &
cd $BANK_SIMULATOR_DIR && uvicorn main:app &
cd $PAYMENT_API_SERVICE_DIR && ./app &
cd $MERCHANT_UI_DIR && nohup npm start >/dev/null 2>&1 &

# cleanup on exit
function handle_interrupt {
    echo "Script interrupted. Exiting..."

    lsof -i ":${MERCHANT_SERVICE_PORT}" | grep -E ^app | awk '{print $2}' | xargs -I@ kill -9 @
    lsof -i ":${RATE_LIMITER_SERVICE_PORT}" | grep -E ^app | awk '{print $2}' | xargs -I@ kill -9 @
    lsof -i ":${LEDGER_SERVICE_PORT}" | grep -E ^app | awk '{print $2}' | xargs -I@ kill -9 @
    lsof -i ":${BANK_SIMULATOR_PORT}" | grep -E ^app | awk '{print $2}' | xargs -I@ kill -9 @
    lsof -i ":${PAYMENT_API_SERVICE_PORT}" | grep -E ^app | awk '{print $2}' | xargs -I@ kill -9 @

    exit 0
}

# Trap the interrupt signal (Ctrl+C) and call the handler function
trap handle_interrupt SIGINT

echo "Press Ctrl+C to exit"

# Infinite loop to keep the script running until interrupted
while true; do
    sleep 1
done