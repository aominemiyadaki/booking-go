package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainConstantsAndErrors(t *testing.T) {
	assert.Equal(t, "отменен", StatusCancelled)
	assert.Equal(t, "Europe/Moscow", DefaultTimezone)
	assert.Equal(t, 8, MinPasswordLength)

	assert.Equal(t, "запись не найдена", ErrNotFound.Error())
	assert.Equal(t, "выбранное время уже занято", ErrSlotTaken.Error())
}
