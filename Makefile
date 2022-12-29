
GOCMD?= go

.PHONY: otelcol
otelcol:
	mkdir -p bin
	cd ./cmd/otto && GO111MODULE=on CGO_ENABLED=0 $(GOCMD) build -trimpath -o ../../bin/otelcol  .