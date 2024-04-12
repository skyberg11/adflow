package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token1, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	intruder, err := client.createUser("Andrew", "Ivanov", "abacaba", "12345678", "arr@mail.ru", "+79821233123")
	assert.NoError(t, err)

	token2, err := client.loginUser("abacaba", "12345678")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world", token1.Token)
	assert.NoError(t, err)

	_, err = client.changeAdStatus(intruder.Data.ID, resp.Data.ID, true, token2.Token)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token1, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	intruder, err := client.createUser("Andrew", "Ivanov", "abacaba", "12345678", "arr@mail.ru", "+79821233123")
	assert.NoError(t, err)

	token2, err := client.loginUser("abacaba", "12345678")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world", token1.Token)
	assert.NoError(t, err)

	_, err = client.updateAd(intruder.Data.ID, resp.Data.ID, "title", "text", token2.Token)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token1, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world", token1.Token)
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(user.Data.ID, "hello", "world", token1.Token)
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))

	resp, err = client.createAd(user.Data.ID, "hello", "world", token1.Token)
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(3))
}
