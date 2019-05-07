SRC = $(shell find . -name "*.go")

run: concourse_db_exporter
	./$< \
		--postgres.database=concourse \
		--postgres.host=127.0.0.1 \
		--postgres.port=6543 \
		--postgres.user=dev \
		--postgres.password=dev

concourse_db_exporter: $(SRC)
	go build -o $@ -v -i .
