test: memcached_stop memcached_start
	@echo "Starting tests.."
	@go test $(glide novendor) -v -cover
	@make -s memcached_stop

memcached_start:
	@echo "Starting memcached 1.."
	@memcached -p 21210 -d
	@echo "Starting memcached 2.."
	@memcached -p 21211 -d
	@echo "Starting memcached 3.."
	@memcached -p 21212 -d
	@echo "Starting memcached 4.."
	@memcached -p 21213 -d
	@echo "Starting memcached 5.."
	@memcached -p 21214 -d
	@echo "Starting memcached 6.."
	@memcached -p 21215 -d
	@echo "Starting memcached 7.."
	@memcached -p 21216 -d
	@echo "Starting memcached 8.."
	@memcached -p 21217 -d
	@echo "Starting memcached 9.."
	@memcached -p 21218 -d
	@echo "Starting memcached 10.."
	@memcached -p 21219 -d
	@echo "Done."

memcached_stop:
	@killall -9 memcached 2> /dev/null; true

