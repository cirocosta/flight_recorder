SRC = $(shell find . -name "*.go")

run: flight_recorder
	./$< \
		--postgres-database=concourse \
		--postgres-host=127.0.0.1 \
		--postgres-port=6543 \
		--postgres-user=dev \
		--postgres-password=dev

flight_recorder: $(SRC)
	go build -o $@ -v -i .
