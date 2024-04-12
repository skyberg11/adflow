package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadToken(t *testing.T) {
	client := getTestClient()
	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "hello", "world", "42istheanswerforanyquestion")
	assert.ErrorIs(t, err, ErrBadRequest)

	_, err = client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)
}

func TestAnotherToken(t *testing.T) {
	client := getTestClient()
	_, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	another, err := client.createUser("Andrew", "Ivanov", "abacaba", "12345678", "arr@mail.ru", "+79821233123")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	_, err = client.createAd(another.Data.ID, "hello", "world", token.Token)
	assert.ErrorIs(t, err, ErrForbidden)
}
