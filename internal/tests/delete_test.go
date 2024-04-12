package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	client := getTestClient()
	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token1, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	another, err := client.createUser("Andrew", "Ivanov", "abacaba", "12345678", "arr@mail.ru", "+79821233123")
	assert.NoError(t, err)

	token2, err := client.loginUser("abacaba", "12345678")
	assert.NoError(t, err)

	_, err = client.createAd(another.Data.ID, "hello", "world", token2.Token)
	assert.NoError(t, err)

	_, err = client.deleteUser(user.Data.ID, another.Data.ID, token2.Token)
	assert.Error(t, err)

	_, err = client.deleteUser(user.Data.ID, user.Data.ID, token1.Token)
	assert.NoError(t, err)

	_, err = client.getUser(user.Data.ID)
	assert.Error(t, err)
}

func TestAdDeletion(t *testing.T) {
	client := getTestClient()
	user1, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token1, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	user2, err := client.createUser("Andrew", "Ivanov", "abacaba", "12345678", "arr@mail.ru", "+79821233123")
	assert.NoError(t, err)

	token2, err := client.loginUser("abacaba", "12345678")
	assert.NoError(t, err)

	ad1Resp, err := client.createAd(user1.Data.ID, "hello", "world", token1.Token)
	assert.NoError(t, err)

	ad2Resp, err := client.createAd(user2.Data.ID, "hello", "world", token2.Token)
	assert.NoError(t, err)

	_, err = client.deleteAd(ad1Resp.Data.ID, user1.Data.ID, token1.Token)
	assert.NoError(t, err)

	_, err = client.getAd(ad1Resp.Data.ID)
	assert.Error(t, err)

	_, err = client.deleteAd(ad2Resp.Data.ID, user1.Data.ID, token1.Token)
	assert.Error(t, err)
}
