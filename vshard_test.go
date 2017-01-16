package vshard

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type VShardTestSuite struct {
	suite.Suite
	Pool *Pool
}

func (suite *VShardTestSuite) SetupSuite() {
	suite.Pool = setupPool(suite.T())
}

func (suite *VShardTestSuite) TearDownTest() {
	tearDownPool(suite.T(), suite.Pool)
}

func (suite *VShardTestSuite) TestNumberOfServers() {
	suite.Equal(10, suite.Pool.numServers)
	suite.Len(suite.Pool.Servers, 10)
	suite.Len(suite.Pool.pool, 10)
}

func (suite *VShardTestSuite) TestStatus() {
	capacity := int64(10)
	available := int64(10)
	maxCap := int64(10)
	idleTimeout := time.Second * 5

	statusList := suite.Pool.Status()
	suite.Len(statusList, 10)

	for _, status := range statusList {
		suite.Equal(capacity, status.Capacity)
		suite.Equal(available, status.Available)
		suite.Equal(maxCap, status.MaxCap)
		suite.Equal(idleTimeout, status.IdleTimeout)
	}
}

func (suite *VShardTestSuite) testMD5Sharding(key string, poolNum int) {
	actualPoolNum := ShardedServerStrategyMD5(key, 10)
	suite.Equal(poolNum, actualPoolNum)
}

func (suite *VShardTestSuite) testFarmhashSharding(key string, poolNum int) {
	actualPoolNum := ShardedServerStrategyFarmhash(key, 10)
	suite.Equal(poolNum, actualPoolNum)
}

func (suite *VShardTestSuite) testShardingDistribution(key, value string, poolNum int) {
	ok, err := suite.Pool.Set(key, 0, 0, []byte(value))

	suite.True(ok)
	suite.NoError(err)

	resource, err := suite.Pool.GetPoolConnection(poolNum)
	defer suite.Pool.ReturnConnection(poolNum, resource)
	if err != nil {
		suite.FailNow("Failure getting specific connection from pool", err)
	}

	result, err := resource.Get(key)
	if err != nil {
		suite.FailNow("Failure getting key", err, result)
	}

	if suite.NotEmpty(result) {
		suite.Equal(value, string(result[0].Value))
	}
}

func (suite *VShardTestSuite) TestShardingDistributionMD5() {
	oldServerStrategy := suite.Pool.ServerStrategy
	suite.Pool.ServerStrategy = ShardedServerStrategyMD5

	suite.testShardingDistribution("f", "test-server-1", 0)
	suite.testShardingDistribution("o", "test-server-2", 1)
	suite.testShardingDistribution("d", "test-server-3", 2)
	suite.testShardingDistribution("e", "test-server-4", 3)
	suite.testShardingDistribution("b", "test-server-5", 4)
	suite.testShardingDistribution("g", "test-server-6", 5)
	suite.testShardingDistribution("p", "test-server-7", 6)
	suite.testShardingDistribution("c", "test-server-8", 7)
	suite.testShardingDistribution("l", "test-server-9", 8)
	suite.testShardingDistribution("a", "test-server-10", 9)

	suite.Pool.ServerStrategy = oldServerStrategy
}

func (suite *VShardTestSuite) TestShardingDistributionFarmhash() {
	oldServerStrategy := suite.Pool.ServerStrategy
	suite.Pool.ServerStrategy = ShardedServerStrategyFarmhash

	suite.testShardingDistribution("f", "test-server-1", 0)
	suite.testShardingDistribution("m", "test-server-2", 1)
	suite.testShardingDistribution("5", "test-server-3", 2)
	suite.testShardingDistribution("l", "test-server-4", 3)
	suite.testShardingDistribution("h", "test-server-5", 4)
	suite.testShardingDistribution("d", "test-server-6", 5)
	suite.testShardingDistribution("a", "test-server-7", 6)
	suite.testShardingDistribution("za", "test-server-8", 7)
	suite.testShardingDistribution("i", "test-server-9", 8)
	suite.testShardingDistribution("b", "test-server-10", 9)

	suite.Pool.ServerStrategy = oldServerStrategy
}

func (suite *VShardTestSuite) TestMD5Sharding() {
	suite.testMD5Sharding("f", 0)
	suite.testMD5Sharding("o", 1)
	suite.testMD5Sharding("d", 2)
	suite.testMD5Sharding("e", 3)
	suite.testMD5Sharding("b", 4)
	suite.testMD5Sharding("g", 5)
	suite.testMD5Sharding("p", 6)
	suite.testMD5Sharding("c", 7)
	suite.testMD5Sharding("l", 8)
	suite.testMD5Sharding("a", 9)
}

func (suite *VShardTestSuite) TestFarmhashSharding() {
	suite.testFarmhashSharding("f", 0)
	suite.testFarmhashSharding("m", 1)
	suite.testFarmhashSharding("5", 2)
	suite.testFarmhashSharding("l", 3)
	suite.testFarmhashSharding("h", 4)
	suite.testFarmhashSharding("d", 5)
	suite.testFarmhashSharding("a", 6)
	suite.testFarmhashSharding("za", 7)
	suite.testFarmhashSharding("i", 8)
	suite.testFarmhashSharding("b", 9)
}

func TestVShardTestSuite(t *testing.T) {
	suite.Run(t, new(VShardTestSuite))
}

// from fib_test.go
func BenchmarkGetKeyMappingMD5(b *testing.B) {
	servers := []string{"0"}
	pool := Pool{
		Servers:        servers,
		ServerStrategy: ShardedServerStrategyMD5,
		numServers:     len(servers),
	}
	keys := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_ = pool.GetKeyMapping(keys...)
	}
}

func BenchmarkGetKeyMappingFarmhash(b *testing.B) {
	servers := []string{"0"}
	pool := Pool{
		Servers:        servers,
		ServerStrategy: ShardedServerStrategyFarmhash,
		numServers:     len(servers),
	}
	keys := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_ = pool.GetKeyMapping(keys...)
	}
}

func BenchmarkShardedServerStrategyMD5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = ShardedServerStrategyMD5("a", 10)
	}
}

func BenchmarkShardedServerStrategyFarmHash(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = ShardedServerStrategyFarmhash("a", 10)
	}
}
