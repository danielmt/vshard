package vshard

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestServers() []string {
	servers := []string{
		"127.0.0.1:21210",
		"127.0.0.1:21211",
		"127.0.0.1:21212",
		"127.0.0.1:21213",
		"127.0.0.1:21214",
		"127.0.0.1:21215",
		"127.0.0.1:21216",
		"127.0.0.1:21217",
		"127.0.0.1:21218",
		"127.0.0.1:21219",
	}

	return servers
}

func setupPool(t assert.TestingT) *Pool {
	pool := Pool{
		Servers:     getTestServers(),
		Capacity:    10,
		MaxCapacity: 10,
		IdleTimeout: time.Second * 5,
	}
	pool.Start()

	return &pool
}

func tearDownPool(t assert.TestingT, pool *Pool) {
	errs := pool.FlushAll()
	for _, err := range errs {
		if err != nil {
			assert.FailNow(t, "Failure on FlushAll", err)
		}
	}
}
