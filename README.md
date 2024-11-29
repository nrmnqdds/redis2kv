redis2kv
========

A simple tool to migrate all keys from a Redis instance to a Cloudflare KV store.

TUTORIAL
--------

1.	Git clone this repository

	```bash
	git clone https://github.com/nrmnqdds/redis2kv.git
	```

2.	Install the required packages

	```bash
	go mod tidy
	```

3.	Copy the `.env.example` file to `.env` and fill in the required values

	```bash
	cp .env.example .env
	```

4.	Run the tool

	```bash
	go run main.go
	```

5.	Profit!
