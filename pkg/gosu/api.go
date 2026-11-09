package gosu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// maxBulkIDs is the id ceiling osu! enforces on its bulk users and beatmaps
// endpoints. A request over it is rejected before any call is made.
const maxBulkIDs = 50

// ErrTooManyIDs is returned by the bulk fetches when more ids are passed than osu!
// accepts in one request.
var ErrTooManyIDs = fmt.Errorf("at most %d ids per bulk request", maxBulkIDs)

// bulkIDsURL builds an endpoint that repeats the ids[] query parameter, the form
// osu!'s bulk users and beatmaps endpoints read a list of ids from.
func bulkIDsURL(endpoint string, ids []int64) string {
	q := url.Values{}
	for _, id := range ids {
		q.Add("ids[]", strconv.FormatInt(id, 10))
	}
	return APIv2URL(endpoint + "?" + q.Encode())
}

type Credentials struct {
	ClientID     int
	ClientSecret string
}

// APIVersion is the osu! API v2 version this package requests, sent as the
// x-api-version header on every v2 call. osu-web gates response-shape changes on
// it and compares it as an integer, so pinning the newest published version keeps
// responses on the current solo_score shape. Update this when osu! ships a newer
// version whose shape this package has been adjusted to decode.
const APIVersion = "20241024"

// apiVersionHeader is the header name osu-web reads the version from.
const apiVersionHeader = "x-api-version"

func buildOAUTHUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", OAuthURL, endpoint)
}

func APIv2URL(endpoint string) string {
	return fmt.Sprintf("%s/%s", ApiV2, endpoint)
}

// GetUserID resolves a username to its user id or ErrUserNotFound
func (c *Client) GetUserID(username string) (int64, error) {
	user, err := c.GetUserByName(username)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func decodeUser(res *http.Response, user any) error {
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return ErrUserNotFound
	}
	if res.StatusCode != 200 {
		return &StatusError{Code: res.StatusCode, Status: res.Status}
	}

	return json.NewDecoder(res.Body).Decode(user)
}

// GetUser fetches a profile by user id or ErrUserNotFound
func (c *Client) GetUser(id int64) (*UserExtended, error) {
	req, _ := http.NewRequest(http.MethodGet, APIv2URL(fmt.Sprintf("users/%d/osu", id)), nil)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	user := &UserExtended{}
	if err := decodeUser(res, user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByName fetches a profile by username or ErrUserNotFound
func (c *Client) GetUserByName(username string) (*UserExtended, error) {
	req, _ := http.NewRequest(http.MethodGet, APIv2URL(fmt.Sprintf("users/%s/osu?key=username", url.PathEscape(username))), nil)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	user := &UserExtended{}
	if err := decodeUser(res, user); err != nil {
		return nil, err
	}
	return user, nil
}

// getUserScores fetches one page of users/{id}/scores/{kind}, decoding into the
// score shape that kind returns.
func getUserScores[T any](c *Client, userID int64, kind string, limit, offset int) ([]T, error) {
	endpoint := fmt.Sprintf("users/%d/scores/%s?mode=osu&limit=%d&offset=%d", userID, kind, limit, offset)

	req, _ := http.NewRequest(http.MethodGet, APIv2URL(endpoint), nil)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrUserNotFound
	}
	if res.StatusCode != 200 {
		return nil, &StatusError{Code: res.StatusCode, Status: res.Status}
	}

	scores := make([]T, 0)
	if err := json.NewDecoder(res.Body).Decode(&scores); err != nil {
		return nil, err
	}

	return scores, nil
}

// GetRecentScores fetches one page of a user's recent scores newest first or ErrUserNotFound
func (c *Client) GetRecentScores(userID int64, limit, offset int) (FullScores, error) {
	return getUserScores[FullScore](c, userID, "recent", limit, offset)
}

// GetBestScores fetches one page of a user's top plays ordered by pp descending or ErrUserNotFound
// osu! clamps limit between 1 and 100 and keeps at most 200 best scores per player
func (c *Client) GetBestScores(userID int64, limit, offset int) (BestScores, error) {
	return getUserScores[BestScore](c, userID, "best", limit, offset)
}

// GetScores fetches one page of the ruleset's global scores feed by cursor and returns the next cursor
// Scores carry a beatmap_id but no embedded beatmap
func (c *Client) GetScores(ruleset, cursor string) (Scores, string, error) {
	endpoint := "scores?ruleset=" + ruleset
	if cursor != "" {
		endpoint += "&cursor_string=" + cursor
	}

	req, _ := http.NewRequest(http.MethodGet, APIv2URL(endpoint), nil)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, "", &StatusError{Code: res.StatusCode, Status: res.Status}
	}

	page := struct {
		Scores       Scores `json:"scores"`
		CursorString string `json:"cursor_string"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, "", err
	}
	return page.Scores, page.CursorString, nil
}

// GetBeatmapScores fetches all of a user's scores on one beatmap
func (c *Client) GetBeatmapScores(userID int64, beatmapID int64) (Scores, error) {
	endpoint := fmt.Sprintf("beatmaps/%d/scores/users/%d/all", beatmapID, userID)

	req, _ := http.NewRequest(http.MethodGet, APIv2URL(endpoint), nil)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, &StatusError{Code: res.StatusCode, Status: res.Status}
	}

	s := struct {
		Scores Scores `json:"scores"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		return nil, err
	}

	return s.Scores, nil
}

// GetBeatmap fetches a beatmap and its owning beatmapset by id
func (c *Client) GetBeatmap(id int64) (FullBeatmap, error) {
	req, _ := http.NewRequest(http.MethodGet, APIv2URL(fmt.Sprintf("beatmaps/%d", id)), nil)
	res, err := c.http.Do(req)
	if err != nil {
		return FullBeatmap{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return FullBeatmap{}, &StatusError{Code: res.StatusCode, Status: res.Status}
	}

	body := FullBeatmap{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return FullBeatmap{}, err
	}
	return body, nil
}

// GetBeatmapAttributes fetches a beatmap's difficulty under a mod combination in its own ruleset
// Mods are acronyms
func (c *Client) GetBeatmapAttributes(id int64, mods []string) (DifficultyAttributes, error) {
	payload, err := json.Marshal(struct {
		Mods []string `json:"mods,omitempty"`
	}{Mods: mods})
	if err != nil {
		return DifficultyAttributes{}, err
	}

	req, _ := http.NewRequest(http.MethodPost, APIv2URL(fmt.Sprintf("beatmaps/%d/attributes", id)), bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res, err := c.http.Do(req)
	if err != nil {
		return DifficultyAttributes{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return DifficultyAttributes{}, &StatusError{Code: res.StatusCode, Status: res.Status}
	}

	body := struct {
		Attributes DifficultyAttributes `json:"attributes"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return DifficultyAttributes{}, err
	}
	return body.Attributes, nil
}

// GetUsers fetches profiles for many ids in one request keyed by id in any order
// Returns ErrTooManyIDs above maxBulkIDs
func (c *Client) GetUsers(ids []int64) ([]*UserWithRulesets, error) {
	if len(ids) > maxBulkIDs {
		return nil, ErrTooManyIDs
	}
	if len(ids) == 0 {
		return nil, nil
	}

	req, _ := http.NewRequest(http.MethodGet, bulkIDsURL("users", ids), nil)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, &StatusError{Code: res.StatusCode, Status: res.Status}
	}

	page := struct {
		Users []*UserWithRulesets `json:"users"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, err
	}
	return page.Users, nil
}

// GetBeatmaps fetches metadata for many beatmap ids in one request keyed by id in any order
// Returns ErrTooManyIDs above maxBulkIDs
func (c *Client) GetBeatmaps(ids []int64) ([]FullBeatmap, error) {
	if len(ids) > maxBulkIDs {
		return nil, ErrTooManyIDs
	}
	if len(ids) == 0 {
		return nil, nil
	}

	req, _ := http.NewRequest(http.MethodGet, bulkIDsURL("beatmaps", ids), nil)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, &StatusError{Code: res.StatusCode, Status: res.Status}
	}

	page := struct {
		Beatmaps []FullBeatmap `json:"beatmaps"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, err
	}
	return page.Beatmaps, nil
}
