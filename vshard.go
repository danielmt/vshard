package vshard

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"

	farm "github.com/dgryski/go-farm"
	jump "github.com/dgryski/go-jump"
	"github.com/youtube/vitess/go/memcache"
	"github.com/youtube/vitess/go/pools"
)

var (
	// ErrKeyNotFound defines the error mensage when key is not found on memcached
	ErrKeyNotFound = errors.New("error: key not found")

	defaultServerStrategy  = FarmhashShardServerStrategy
	defaultHashKeyStrategy = FarmhashKeyStrategy
)

const (
	defaultCapacity          = 5
	defaultMaxCapacity       = 5
	defaultIdleTimeout       = time.Millisecond * 500
	defaultConnectionTimeout = time.Millisecond * 200
)

// VitessResource implements the expected interface for vitess internal pool
type VitessResource struct {
	*memcache.Connection
}

// ServerStrategy defines the signature for the sharding function
type ServerStrategy func(key string, numServers int) int

// HashKeyStrategy defines the signature for the key hashing function
type HashKeyStrategy func(key string) string

// Close closes connections in a pool
func (r VitessResource) Close() {
	r.Connection.Close()
}

// Pool defines the pool
type Pool struct {
	Servers               []string
	pool                  []*pools.ResourcePool
	Capacity, MaxCapacity int
	numServers            int
	ServerStrategy        ServerStrategy
	HashKeyStrategy       HashKeyStrategy
	IdleTimeout           time.Duration
	ConnectionTimeout     time.Duration
	sync.RWMutex
}

// PoolStats defines all stats vitess memcached driver exposes
type PoolStats struct {
	Slot        int
	Server      string
	Capacity    int64
	Available   int64
	MaxCap      int64
	WaitCount   int64
	WaitTime    time.Duration
	IdleTimeout time.Duration
}

// MD5ShardServerStrategy uses md5+jump to pick a server
func MD5ShardServerStrategy(key string, numServers int) int {
	if numServers == 1 {
		return 0
	}

	hash := md5.Sum([]byte(key))
	hashInt := big.NewInt(0)
	hashInt.SetString(hex.EncodeToString(hash[:]), 16)

	return int(jump.Hash(hashInt.Uint64(), numServers))
}

// FarmhashShardServerStrategy uses farmhash+jump to pick a server
func FarmhashShardServerStrategy(key string, numServers int) int {
	if numServers == 1 {
		return 0
	}

	return int(jump.Hash(farm.Fingerprint64([]byte(key)), numServers))
}

// FarmhashKeyStrategy uses farmhash to normalize key names for storage
func FarmhashKeyStrategy(key string) string {
	return strconv.FormatUint(farm.Fingerprint64([]byte(key)), 10)
}

// NoKeyStrategy doesn't hash the key
func NoKeyStrategy(key string) string {
	return key
}

// Start starts the pool
func (v *Pool) Start() {

	v.initialize()

	for i, server := range v.Servers {
		func(_v *Pool, _server string, _i int) {
			_v.pool = append(_v.pool, pools.NewResourcePool(func() (pools.Resource, error) {
				c, err := memcache.Connect(_server, _v.ConnectionTimeout)
				return VitessResource{c}, err
			}, _v.Capacity, _v.MaxCapacity, _v.IdleTimeout))

			conn, err := _v.GetPoolConnection(_i)
			defer _v.ReturnConnection(_i, conn)
			if err != nil {
				log.Fatalf("Can't connect to memcached %d (%s): %s", _i, _server, err)
			}
		}(v, server, i)
	}
}

func (v *Pool) initialize() {
	v.numServers = len(v.Servers)
	v.pool = []*pools.ResourcePool{}

	if v.ServerStrategy == nil {
		v.ServerStrategy = defaultServerStrategy
	}
	if v.HashKeyStrategy == nil {
		v.HashKeyStrategy = defaultHashKeyStrategy
	}
	if v.Capacity == 0 {
		v.Capacity = defaultCapacity
	}
	if v.MaxCapacity == 0 {
		v.MaxCapacity = defaultMaxCapacity
	}
	if v.IdleTimeout == 0 {
		v.IdleTimeout = defaultIdleTimeout
	}
	if v.ConnectionTimeout == 0 {
		v.ConnectionTimeout = defaultConnectionTimeout
	}
}

// GetConnection returns a connection from the sharding pool, based on the key
func (v *Pool) GetConnection(key string) (*VitessResource, int, error) {
	poolNum := v.ServerStrategy(key, v.numServers)

	connection, err := v.GetPoolConnection(poolNum)
	if err != nil {
		return nil, -1, err
	}

	return connection, poolNum, nil
}

// GetPoolConnection returns a connection from a specific pool number
func (v *Pool) GetPoolConnection(poolNum int) (*VitessResource, error) {
	ctx := context.Background()

	v.RLock()
	resource, err := v.pool[poolNum].Get(ctx)
	v.RUnlock()
	if err != nil {
		return nil, err
	}

	connection := resource.(VitessResource)

	return &connection, nil
}

// ReturnConnection returns a connection to the pool
func (v *Pool) ReturnConnection(poolNum int, resource *VitessResource) {
	if poolNum > v.numServers || poolNum < 0 {
		log.Fatalf("error: invalid server %d (of total %d)", poolNum, v.numServers)
	}

	v.RLock()
	v.pool[poolNum].Put(*resource)
	v.RUnlock()
}

// GetKeyMapping returns a mapping of server to a list of keys, useful for Gets()
func (v *Pool) GetKeyMapping(keys ...string) map[int][]string {
	mapping := make(map[int][]string)

	for i := 0; i < v.numServers; i++ {
		mapping[i] = []string{}
	}

	for _, key := range keys {
		poolNum := v.ServerStrategy(key, v.numServers)
		mapping[poolNum] = append(mapping[poolNum], v.HashKeyStrategy(key))
	}

	return mapping
}
