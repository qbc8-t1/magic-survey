package test

import (
	"github.com/QBC8-Team1/magic-survey/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	key := "And our minds believed that life was 514v3ry..."
	lifePass, _ := utils.HashPassword(key)

	assert.Nil(t, utils.CheckPasswordHash(key, lifePass))
}
