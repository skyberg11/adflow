package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()
	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(user.Data.ID))
	assert.False(t, response.Data.Published)
}

type Ad struct {
	Title string
	Text  string
}

type TestAd struct {
	In         Ad
	ExpectName string
}

func TestGetAd(t *testing.T) {
	client := getTestClient()

	tests := []TestAd{
		{Ad{"pen", "cool pen"}, "pen"},
		{Ad{"abac", "kek@gmail.com"}, "abac"},
		{Ad{"pc", "write to 123abacaba@bk.ru"}, "pc"},
		{Ad{"test", "1"}, "test"},
	}

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	for idx, test := range tests {
		test := test
		idx := idx

		t.Run(test.ExpectName, func(t *testing.T) {
			t.Parallel()

			response, err := client.createAd(user.Data.ID, test.In.Title, test.In.Text, token.Token)
			assert.NoError(t, err)

			got, err := client.getAd(response.Data.ID)
			assert.NoError(t, err)

			if got.Data.Title != test.ExpectName {
				t.Errorf(`test %d: expect %v got %v`, idx, test.ExpectName, got.Data.Title)
			}
		})
	}
}

func BenchmarkCreateDeleteAd(b *testing.B) {
	client := getTestClient()
	user, _ := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")

	token, _ := client.loginUser("skyberg11", "abacaba")

	for i := 0; i < b.N; i++ {
		response, _ := client.createAd(user.Data.ID, "hello", "world", token.Token)
		_, _ = client.deleteAd(user.Data.ID, response.Data.ID, token.Token)
	}
}

func BenchmarkChangeAdStatus(b *testing.B) {
	client := getTestClient()
	user, _ := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")

	token, _ := client.loginUser("skyberg11", "abacaba")

	response, _ := client.createAd(user.Data.ID, "hello", "world", token.Token)

	for i := 0; i < b.N; i++ {
		_, _ = client.changeAdStatus(user.Data.ID, response.Data.ID, i%2 == 0, token.Token)
	}

	_, _ = client.deleteAd(response.Data.ID, user.Data.ID, token.Token)
}

func TestTimeTest(t *testing.T) {
	client := getTestClient()

	start := time.Now().UTC()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)
	assert.Less(t, start, response.Data.CreationTime)

	response, err = client.updateAd(user.Data.ID, response.Data.ID, "привет", "мир", token.Token)
	assert.NoError(t, err)

	end := time.Now().UTC()

	assert.True(t, start.Before(response.Data.CreationTime))
	assert.True(t, response.Data.CreationTime.Before(response.Data.UpdateTime))
	assert.True(t, response.Data.UpdateTime.Before(end))
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()
	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.AuthorID, int64(user.Data.ID))

	response, err = client.changeAdStatus(user.Data.ID, response.Data.ID, true, token.Token)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)
	assert.Equal(t, response.Data.AuthorID, int64(user.Data.ID))

	response, err = client.changeAdStatus(user.Data.ID, response.Data.ID, false, token.Token)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
	assert.Equal(t, response.Data.AuthorID, int64(user.Data.ID))

	response, err = client.changeAdStatus(user.Data.ID, response.Data.ID, false, token.Token)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
	assert.Equal(t, response.Data.AuthorID, int64(user.Data.ID))
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)

	response, err = client.updateAd(user.Data.ID, response.Data.ID, "привет", "мир", token.Token)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
	assert.Equal(t, response.Data.AuthorID, int64(user.Data.ID))
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("Timur", "Zykov", "skyberg11", "abacaba", "zykov.ta@phystech.edu", "891428821XX")
	assert.NoError(t, err)

	token, err := client.loginUser("skyberg11", "abacaba")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world", token.Token)
	assert.NoError(t, err)
	publishedAd, err := client.changeAdStatus(user.Data.ID, response.Data.ID, true, token.Token)
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "best cat", "not for sale", token.Token)
	assert.NoError(t, err)

	ads, err := client.listAds(1, nil, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}
