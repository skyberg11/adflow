package tests

import (
	"github.com/stretchr/testify/assert"

	"strings"
	"testing"
)

func TestListPublishedAuthor(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token1, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	another, err := client.createUser("Andrew", "Ivanov", "abacaba", "12345678", "arr@mail.ru", "+79821233123")
	assert.NoError(t, err)

	token2, err := client.loginUser("abacaba", "12345678")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world", token1.Token)
	assert.NoError(t, err)

	publishedAd1, err := client.changeAdStatus(user.Data.ID, response.Data.ID, true, token1.Token)
	assert.NoError(t, err)

	response, err = client.createAd(another.Data.ID, "best cat", "not for sale", token2.Token)
	assert.NoError(t, err)

	_, err = client.changeAdStatus(another.Data.ID, response.Data.ID, true, token2.Token)
	assert.NoError(t, err)

	ads, err := client.listAds(1, user.Data.ID, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd1.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd1.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd1.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd1.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestListTitlePrefix(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token1, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world", token1.Token)
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user.Data.ID, response.Data.ID, true, token1.Token)
	assert.NoError(t, err)

	publishedAd, err := client.createAd(user.Data.ID, "best cat", "not for sale", token1.Token)
	assert.NoError(t, err)

	ads, err := client.listAds(nil, nil, "best", nil)

	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
}

func FuzzGenFilter(f *testing.F) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(f, err)

	token1, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(f, err)

	// "hello" "world" true
	// "best cat" "not for sale" false
	f.Add(1, int64(0), "hello")
	f.Add(0, int64(0), "hello")
	f.Add(1, int64(1), "best")
	f.Add(1, int64(1), "bad")

	f.Fuzz(func(t *testing.T, published int, AuthorID int64, TitlePrefix string) {

		response, err := client.createAd(user.Data.ID, "hello", "world", token1.Token)
		assert.NoError(t, err)

		_, err = client.changeAdStatus(user.Data.ID, response.Data.ID, true, token1.Token)
		assert.NoError(t, err)

		_, err = client.createAd(user.Data.ID, "best cat", "not for sale", token1.Token)
		assert.NoError(t, err)

		ads, err := client.listAds(published, AuthorID, TitlePrefix, nil)
		assert.NoError(t, err)

		if published != 0 {
			if AuthorID == user.Data.ID && strings.HasPrefix("hello", TitlePrefix) {
				assert.Len(t, ads.Data, 1)
				assert.Equal(t, ads.Data[0].Title, "hello")
			} else {
				assert.Len(t, ads.Data, 0)
			}
		} else {
			if AuthorID == user.Data.ID && strings.HasPrefix("best cat", TitlePrefix) {
				assert.Len(t, ads.Data, 1)
				assert.Equal(t, ads.Data[0].Title, "best cat")
			} else {
				assert.Len(t, ads.Data, 0)
			}
		}
	})
}
