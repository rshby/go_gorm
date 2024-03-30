package mock

import (
	"github.com/stretchr/testify/mock"
	"go_gorm/config"
)

type ConfigMock struct {
	Mock *mock.Mock
}

func NewConfigMock() *ConfigMock {
	return &ConfigMock{Mock: &mock.Mock{}}
}

func (c *ConfigMock) GetConfig() *config.Config {
	args := c.Mock.Called()
	return args.Get(0).(*config.Config)
}
