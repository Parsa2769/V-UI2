#!/bin/bash

set -e

echo "Fixing go.sum..."
cd backend
go mod tidy
cd ..

echo "Done."
