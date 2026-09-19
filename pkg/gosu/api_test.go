package gosu

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// roundTripFunc serves a canned response for any request, recording the call count, the
// last URL, and the last request headers so tests can assert what was requested.
type roundTripFunc struct {
	body       string
	status     int
	calls      int
	lastURL    string
	lastHeader http.Header
}

func (f *roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	f.calls++
	f.lastURL = req.URL.String()
	f.lastHeader = req.Header.Clone()
	status := f.status
	if status == 0 {
		status = 200
	}
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(f.body)),
		Header:     make(http.Header),
	}, nil
}

// failRoundTrip fails the test if any request is made through it.
type failRoundTrip struct{ t *testing.T }

func (f failRoundTrip) RoundTrip(*http.Request) (*http.Response, error) {
	f.t.Fatal("unexpected API request")
	return nil, errors.New("unreachable")
}

func clientWith(rt http.RoundTripper) *Client {
	return &Client{http: &http.Client{Transport: rt}}
}

func resourceClientWith(rt http.RoundTripper) *ResourceClient {
	return &ResourceClient{Client: clientWith(rt)}
}

func manyIDs(n int) []int64 {
	ids := make([]int64, n)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	return ids
}

func TestBulkIDsURLRepeatsParam(t *testing.T) {
	got := bulkIDsURL("users", []int64{7, 42})
	if !strings.Contains(got, "ids%5B%5D=7") || !strings.Contains(got, "ids%5B%5D=42") {
		t.Fatalf("url %q missing repeated ids[] params", got)
	}
}

func TestGetUsersRejectsTooManyIDs(t *testing.T) {
	c := clientWith(failRoundTrip{t})
	if _, err := c.GetUsers(manyIDs(maxBulkIDs + 1)); !errors.Is(err, ErrTooManyIDs) {
		t.Fatalf("err = %v, want ErrTooManyIDs", err)
	}
}

func TestGetUsersEmptySkipsRequest(t *testing.T) {
	c := clientWith(failRoundTrip{t})
	got, err := c.GetUsers(nil)
	if err != nil || got != nil {
		t.Fatalf("GetUsers(nil) = %v, %v; want nil, nil", got, err)
	}
}

// The bulk users endpoint reports ranking stats per ruleset; GetUsers folds the osu!
// ruleset stats onto Statistics so the shape matches the single-user endpoint.
func TestGetUsersFoldsRulesetStats(t *testing.T) {
	rt := &roundTripFunc{body: `{"users":[{"id":9,"username":"a","statistics_rulesets":{"osu":{"pp":1234.5,"global_rank":10}}}]}`}
	got, err := clientWith(rt).GetUsers([]int64{9})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d profiles, want 1", len(got))
	}
	if got[0].Statistics.PP != 1234.5 {
		t.Fatalf("PP = %v, want 1234.5", got[0].Statistics.PP)
	}
	if got[0].GlobalRank() == nil || *got[0].GlobalRank() != 10 {
		t.Fatalf("GlobalRank = %v, want 10", got[0].GlobalRank())
	}
}

func TestGetBeatmapsRejectsTooManyIDs(t *testing.T) {
	c := clientWith(failRoundTrip{t})
	if _, err := c.GetBeatmaps(manyIDs(maxBulkIDs + 1)); !errors.Is(err, ErrTooManyIDs) {
		t.Fatalf("err = %v, want ErrTooManyIDs", err)
	}
}

func TestGetBeatmapsEmptySkipsRequest(t *testing.T) {
	c := clientWith(failRoundTrip{t})
	got, err := c.GetBeatmaps(nil)
	if err != nil || got != nil {
		t.Fatalf("GetBeatmaps(nil) = %v, %v; want nil, nil", got, err)
	}
}

func TestGetBeatmapsDecodesSet(t *testing.T) {
	rt := &roundTripFunc{body: `{"beatmaps":[{"id":5,"max_combo":600,"beatmapset":{"id":3,"title":"t"}}]}`}
	got, err := clientWith(rt).GetBeatmaps([]int64{5})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != 5 || got[0].MaxCombo != 600 {
		t.Fatalf("beatmap = %+v, want id 5 max_combo 600", got)
	}
	if got[0].Beatmapset.ID != 3 || got[0].Beatmapset.Title != "t" {
		t.Fatalf("set = %+v, want id 3 title t", got[0].Beatmapset)
	}
}

func TestGetOwnUserHitsMeEndpoint(t *testing.T) {
	rt := &roundTripFunc{body: `{"id":5,"username":"me"}`}
	user, err := resourceClientWith(rt).GetOwnUser()
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != 5 {
		t.Errorf("id = %d, want 5", user.ID)
	}
	if !strings.HasSuffix(rt.lastURL, "me/osu") {
		t.Errorf("requested %s, want me/osu", rt.lastURL)
	}
}

func TestGetRecentScoresHitsRecentPath(t *testing.T) {
	rt := &roundTripFunc{body: `[{"id":9,"beatmap":{"id":3}}]`}
	got, err := clientWith(rt).GetRecentScores(7, 50, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != 9 || got[0].Beatmap.ID != 3 {
		t.Fatalf("scores = %+v, want one score id 9 on beatmap 3", got)
	}
	if !strings.Contains(rt.lastURL, "users/7/scores/recent") {
		t.Errorf("requested %s, want users/7/scores/recent", rt.lastURL)
	}
}

func TestGetBestScoresHitsBestPath(t *testing.T) {
	rt := &roundTripFunc{body: `[]`}
	if _, err := clientWith(rt).GetBestScores(7, 50, 10); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"users/7/scores/best", "mode=osu", "limit=50", "offset=10"} {
		if !strings.Contains(rt.lastURL, want) {
			t.Errorf("requested %s, want it to contain %s", rt.lastURL, want)
		}
	}
}

func TestGetBestScoresDecodesWeight(t *testing.T) {
	rt := &roundTripFunc{body: `[{"id":9,"pp":200,"weight":{"percentage":95,"pp":190}}]`}
	got, err := clientWith(rt).GetBestScores(7, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("scores = %+v, want one", got)
	}
	if got[0].PP != 200 {
		t.Errorf("pp = %v, want 200", got[0].PP)
	}
	if got[0].Weight.Percentage != 95 || got[0].Weight.PP != 190 {
		t.Errorf("weight = %+v, want percentage 95 pp 190", got[0].Weight)
	}
}

func TestGetUserScoresNotFound(t *testing.T) {
	rt := &roundTripFunc{status: 404}
	if _, err := clientWith(rt).GetBestScores(7, 1, 0); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("err = %v, want ErrUserNotFound", err)
	}
}
