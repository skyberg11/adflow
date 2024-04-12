package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd_EmptyTitle(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "", "world", token.Token)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongTitle(t *testing.T) {
	client := getTestClient()

	title := strings.Repeat("a", 101)

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, title, "world", token.Token)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_EmptyText(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "title", "", token.Token)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongText(t *testing.T) {
	client := getTestClient()

	text := strings.Repeat("a", 501)

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "title", text, token.Token)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyTitle(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)

	_, err = client.updateAd(user.Data.ID, resp.Data.ID, "", "new_world", token.Token)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongTitle(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)

	title := strings.Repeat("a", 101)

	_, err = client.updateAd(user.Data.ID, resp.Data.ID, title, "world", token.Token)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyText(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)

	_, err = client.updateAd(user.Data.ID, resp.Data.ID, "title", "", token.Token)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongText(t *testing.T) {
	client := getTestClient()

	text := strings.Repeat("a", 501)

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	resp, err := client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)

	_, err = client.updateAd(user.Data.ID, resp.Data.ID, "title", text, token.Token)
	assert.ErrorIs(t, err, ErrBadRequest)
}
