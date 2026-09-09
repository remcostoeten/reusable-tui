#!/bin/fish

set entrypoint './cmd/tui'
echo "Running reusable-tui..."
go run $entrypoint $argv
