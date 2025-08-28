#!/bin/bash

echo "🚀 Starting bulk import of participants..."

# Make sure we're in the right directory
cd "$(dirname "$0")/../.."

# Check if participants.json exists
if [ ! -f "participants.json" ]; then
    echo "❌ participants.json not found in current directory"
    exit 1
fi

# Check if data.db exists
if [ ! -f "data.db" ]; then
    echo "❌ data.db not found in current directory"
    exit 1
fi

# Run the bulk import
go run cmd/bulk_import/main.go

echo "✅ Bulk import completed!"
