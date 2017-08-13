export CGO_ENABLED=1

test: memcached_stop memcached_start _test bench memcached_stop

_test:
	@echo "Starting tests.."
	@go test $(glide novendor) -v -cover -race
	@go vet $(glide novendor)

bench:
	@echo "Starting benchmarks.."
	@go test -run=XXX -bench=. -v

memcached_start:
	@for i in `seq 10 19`; do memcached -o modern -p 212$$i -d; done

memcached_stop:
	@killall -9 memcached 2> /dev/null; true

key_dist:
	@go run tools/hash_distribution/main.go
