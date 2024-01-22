package mock

import "github.com/stretchr/testify/mock"

type Validator struct {
	mock.Mock
}

func (m *Validator) Validate(solution string) error {
	return m.Called(solution).Error(0)
}
