package mock

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type Clock struct {
	mock.Mock
}

func (c *Clock) Now() time.Time {
	return c.Called().Get(0).(time.Time)
}
