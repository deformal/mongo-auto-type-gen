default:
    @just --list

run:
    @echo "Local go run"
    @../core/ go run ./core/cmd/main.go --env-file .env --out ./generated
