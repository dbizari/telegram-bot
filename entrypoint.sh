#!/bin/bash
set -e

echo "Building app..."
cd $TELEGRAM_BOT_PATH
make build
echo "App built successfully!"

echo "Running app..."
exec $TELEGRAM_BOT_PATH/build/telegram-bot

exec "$@"