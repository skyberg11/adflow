package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserUniqueSameNick(t *testing.T) {
	client := getTestClient()
	_, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	_, err = client.createUser("Timur1", "Zykov1", "skyberg11", "abacaba", "zykov.ta1@phystech.edu", "891428821XX")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()
	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	doppelganger, err := client.updateUser("Timur", "Zykov", "zykov.t1a@phystech.edu", "891428821XX", user.Data.ID, token.Token)
	assert.NoError(t, err)

	assert.Equal(t, doppelganger.Data.Nickname, user.Data.Nickname)
	assert.Equal(t, doppelganger.Data.ID, user.Data.ID)
}
