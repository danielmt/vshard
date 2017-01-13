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
	suite.Equal(expectedValue, string(value))
}

func (suite *VShardCommandsTestSuite) TestReplaceInexistentKey() {
	key := "replace-key-does-not-exist"
	expectedValue := "vshard-test-replace"
	ok, err := suite.Pool.Replace(key, 0, 0, []byte(expectedValue))
	suite.False(ok)
	suite.NoError(err)
}

func (suite *VShardCommandsTestSuite) TestReplace() {
	key := "replace-key"
	initialValue := "initial-value"
	ok, err := suite.Pool.Set(key, 0, 0, []byte(initialValue))
	suite.True(ok)
	suite.NoError(err)

	expectedValue := "replaced-value"
	ok, err = suite.Pool.Replace(key, 0, 0, []byte(expectedValue))
	suite.True(ok)
	suite.NoError(err)

	value, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(expectedValue, string(value))
}

func (suite *VShardCommandsTestSuite) TestAppend() {
	key := "append-key"
	initialValue := "initial-value"
	ok, err := suite.Pool.Set(key, 0, 0, []byte(initialValue))
	suite.True(ok)
	suite.NoError(err)

	appendValue := ":second-value"
	ok, err = suite.Pool.Append(key, 0, 0, []byte(appendValue))
	suite.True(ok)
	suite.NoError(err)

	value, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(initialValue+appendValue, string(value))
}

func (suite *VShardCommandsTestSuite) TestPrepend() {
	key := "prepend-key"
	initialValue := "initial-value"
	ok, err := suite.Pool.Set(key, 0, 0, []byte(initialValue))
	suite.True(ok)
	suite.NoError(err)

	prependValue := "second-value:"
	ok, err = suite.Pool.Prepend(key, 0, 0, []byte(prependValue))
	suite.True(ok)
	suite.NoError(err)

	value, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(prependValue+initialValue, string(value))
}

func (suite *VShardCommandsTestSuite) TestDelete() {
	key := "delete-key"
	expectedValue := "delete-test"
	ok, err := suite.Pool.Set(key, 0, 0, []byte(expectedValue))
	suite.True(ok)
	suite.NoError(err)

	value, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(expectedValue, string(value))

	ok, err = suite.Pool.Delete(key)
	suite.True(ok)
	suite.NoError(err)

	ok, err = suite.Pool.Delete(key)
	suite.False(ok)
	suite.NoError(err)

	value, err = suite.Pool.Get(key)
	suite.EqualError(err, ErrKeyNotFound.Error())
	suite.Equal("", string(value))
}

func (suite *VShardCommandsTestSuite) TestFlushAll() {
	key := "flush-key"
	expectedValue := "hello-test"
	ok, err := suite.Pool.Set(key, 0, 0, []byte(expectedValue))
	suite.True(ok)
	suite.NoError(err)

	value, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(expectedValue, string(value))

	errs := suite.Pool.FlushAll()
	for _, err := range errs {
		suite.NoError(err)
	}

	value, err = suite.Pool.Get(key)
	suite.EqualError(err, ErrKeyNotFound.Error())
	suite.Equal("", string(value))
}

func (suite *VShardCommandsTestSuite) TestCasFailure() {
	key := "cas-failure-key"
	expectedValue := "set-before-cas-test"
	ok, err := suite.Pool.Set(key, 0, 0, []byte(expectedValue))
	suite.True(ok)
	suite.NoError(err)

	value, err := suite.Pool.Gets(key)
	suite.NoError(err)
	suite.Equal(expectedValue, string(value[0].Value))

	newValue := "updated-value"
	ok, err = suite.Pool.Set(key, 0, 0, []byte(newValue))
	suite.True(ok)
	suite.NoError(err)

	casValue := "cas-value"
	ok, err = suite.Pool.Cas(key, 0, 0, []byte(casValue), value[0].Cas)
	suite.False(ok, "Update should have failed")
	suite.NoError(err)

	afterCasValue, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(newValue, string(afterCasValue), "Should have the second Set() value, not Cas() value")
}

func (suite *VShardCommandsTestSuite) TestCasSuccess() {
	key := "cas-success-key"
	expectedValue := "set-before-cas-test2"
	ok, err := suite.Pool.Set(key, 0, 0, []byte(expectedValue))
	suite.True(ok)
	suite.NoError(err)

	value, err := suite.Pool.Gets(key)
	suite.NoError(err)
	suite.Equal(expectedValue, string(value[0].Value))

	casValue := "cas-value"
	ok, err = suite.Pool.Cas(key, 0, 0, []byte(casValue), value[0].Cas)
	suite.True(ok, "Update should have worked, there was no Set() after Gets()")
	suite.NoError(err)

	afterCasValue, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(casValue, string(afterCasValue), "Should have the Cas() value")
}

func TestVShardCommandsTestSuite(t *testing.T) {
	suite.Run(t, new(VShardCommandsTestSuite))
}
