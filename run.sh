#!/bin/bash

# Build and run Market Data Insights

echo "Building market-data-insights..."
go build -o market-data-insights main.go

if [ $? -eq 0 ]; then
    echo "Build successful. Starting server..."
    ./market-data-insights
else
    echo "Build failed."
    exit 1
fi