package vshard

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type VShardCommandsTestSuite struct {
	suite.Suite
	Pool *Pool
}

func (suite *VShardCommandsTestSuite) SetupSuite() {
	suite.Pool = setupPool(suite.T())
}

func (suite *VShardCommandsTestSuite) TearDownTest() {
	tearDownPool(suite.T(), suite.Pool)
}

func (suite *VShardCommandsTestSuite) TestSetGet() {
	key := "get-set-key"
	expectedValue := "hello-test"
	ok, err := suite.Pool.Set(key, 0, 0, []byte(expectedValue))
	suite.True(ok)
	suite.NoError(err)

	value, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(expectedValue, string(value))
}

func (suite *VShardCommandsTestSuite) TestGetInexistentKey() {
	key := "key-do-not-exist"
	value, err := suite.Pool.Get(key)
	suite.Equal(ErrKeyNotFound, err)
	suite.Empty(value)
}

func (suite *VShardCommandsTestSuite) TestAdd() {
	key := "add-key"
	expectedValue := "vshard-test-add"
	ok, err := suite.Pool.Add(key, 0, 0, []byte(expectedValue))
	suite.True(ok)
	suite.NoError(err)

	newValue := "this-should-not-work"
	ok, err = suite.Pool.Add(key, 0, 0, []byte(newValue))
	suite.False(ok)
	suite.NoError(err)

	value, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(expected, string(value))
}

func TestVShardCommandsTestSuite(t *testing.T) {
	suite.Run(t, new(VShardCommandsTestSuite))
}
