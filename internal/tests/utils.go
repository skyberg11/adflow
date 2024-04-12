package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"adflow/internal/adapters"
	"adflow/internal/adapters/adrepo"
	"adflow/internal/adapters/aduser"
	"adflow/internal/ads"
	"adflow/internal/app"
	"adflow/internal/ports/httpgin"
)

type adData struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type userData struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Nickname   string `json:"nickname"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	ID         int64  `json:"user_id"`
}

type userSecureData struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	ID         int64  `json:"user_id"`
}

type adResponse struct {
	Data adData `json:"data"`
}
type userResponse struct {
	Data userSecureData `json:"data"`
}
type userResponseToken struct {
	Token string `json:"token"`
}
type adsResponse struct {
	Data []adData `json:"data"`
}
type deleteResponse struct {
	Data string `json:"data"`
}

var (
	ErrBadRequest = fmt.Errorf("bad request")
	ErrForbidden  = fmt.Errorf("forbidden")
)

type testClient struct {
	client  *http.Client
	baseURL string
}

func getTestClient() *testClient {
	db_users, err := adapters.NewSQLite("test_users.db")
	db_users.Migrator().DropTable(&ads.User{})

	if err != nil {
		panic(err)
	}

	db_ads, err := adapters.NewSQLite("test_ads.db")
	db_ads.Migrator().DropTable(&ads.Ad{})

	if err != nil {
		panic(err)
	}

	repo, users := adrepo.NewSQLiteAds(db_ads), aduser.NewSQLiteUsers(db_users)

	server := httpgin.NewHTTPServer(":18080", app.NewApp(repo, users))
	testServer := httptest.NewServer(server.Handler())

	return &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func (tc *testClient) getResponse(req *http.Request, out any) error {
	resp, err := tc.client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return ErrBadRequest
		}
		if resp.StatusCode == http.StatusForbidden {
			return ErrForbidden
		}
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("unable to read response: %w", err)
	}

	err = json.Unmarshal(respBody, out)

	if err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}

	return nil
}

func (tc *testClient) getAd(id int64) (adResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", id), nil)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adResponse
	err = tc.getResponse(req, &response)

	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) getUser(id int64) (userResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", id), nil)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response userResponse
	err = tc.getResponse(req, &response)

	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) createAd(userID int64, title string, text string, token string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) createUser(first_name, second_name, nickname, password, email, phone string) (userResponse, error) {
	body := map[string]any{
		"first_name":  first_name,
		"second_name": second_name,
		"nickname":    nickname,
		"password":    password,
		"email":       email,
		"phone":       phone,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/users", bytes.NewReader(data))
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) loginUser(nickname, password string) (userResponseToken, error) {
	body := map[string]any{
		"nickname": nickname,
		"password": password,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponseToken{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/users/login", bytes.NewReader(data))
	if err != nil {
		return userResponseToken{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponseToken
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponseToken{}, err
	}

	return response, nil
}

func (tc *testClient) changeAdStatus(userID int64, adID int64, published bool, token string) (adResponse, error) {
	body := map[string]any{
		"user_id":   userID,
		"published": published,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d/status", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateAd(userID int64, adID int64, title string, text string, token string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) deleteAd(adID int64, userID int64, token string) (deleteResponse, error) {
	body := map[string]any{
		"user_id": userID,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return deleteResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID), bytes.NewReader(data))
	if err != nil {
		return deleteResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	var response deleteResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return deleteResponse{}, err
	}
	return response, err
}

func (tc *testClient) deleteUser(id int64, UserID int64, token string) (deleteResponse, error) {
	body := map[string]any{
		"user_id": id,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return deleteResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", id), bytes.NewReader(data))
	if err != nil {
		return deleteResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	var response deleteResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return deleteResponse{}, err
	}
	return response, err
}

func (tc *testClient) updateUser(first_name, second_name, email, phone string, id int64, token string) (userResponse, error) {
	body := map[string]any{
		"first_name":  first_name,
		"second_name": second_name,
		"email":       email,
		"phone":       phone,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", id), bytes.NewReader(data))
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAds(Published, AuthorID, TitlePrefix, CreationTime any) (adsResponse, error) {
	str := tc.baseURL + "/api/v1/ads"
	pref := "?"
	if Published != nil {
		str = fmt.Sprintf("%s%s%s%d", str, pref, "published=", Published)
		pref = "&"
	}
	if AuthorID != nil {
		str = fmt.Sprintf("%s%s%s%d", str, pref, "author=", AuthorID)
		pref = "&"
	}
	if TitlePrefix != nil {
		str = fmt.Sprintf("%s%s%s%s", str, pref, "title=", TitlePrefix)
		pref = "&"
	}
	if CreationTime != nil {
		str = fmt.Sprintf("%s%s%s%s", str, pref, "creation=", CreationTime)
	}

	req, err := http.NewRequest(http.MethodGet, str, nil)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adsResponse
	err = tc.getResponse(req, &response)

	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}
