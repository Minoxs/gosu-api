package gosu

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

// The fixtures repeat the keys of a live response, every value non-zero so assertFilled
// can tell a mapped field from an unmapped one
const userJSON = `"id":2,"username":"a","avatar_url":"https://a","country_code":"AU","default_group":"default",
	"is_active":true,"is_bot":true,"is_deleted":true,"is_online":true,"is_supporter":true,
	"last_visit":"2026-01-02T03:04:05Z","pm_friends_only":true,"profile_colour":"#3366FF"`

const userAttributesJSON = `"country":{"code":"AU","name":"Australia"},"cover":{"custom_url":"https://c","id":"4","url":"https://c"},
	"groups":[{"colour":"#0066FF","has_listing":true,"has_playmodes":true,"id":33,"identifier":"ppy","is_probationary":true,
	"name":"ppy","playmodes":["osu"],"short_name":"PPY"}],
	"team":{"flag_url":"https://f","id":1,"name":"team","short_name":"T"}`

const userStatisticsJSON = `"accuracy":0.97,"hit_accuracy":97.0,"count_100":1,"count_300":2,"count_50":3,"count_miss":4,
	"global_rank":5,"global_rank_percent":0.25,"grade_counts":{"a":1,"s":2,"sh":3,"ss":4,"ssh":5},"is_ranked":true,
	"level":{"current":60,"progress":55},"maximum_combo":6,"play_count":7,"play_time":8,"pp":9.5,"ranked_score":5000000000,
	"replays_watched_by_others":10,"total_hits":6000000000,"total_score":7000000000`

const userStatisticsExtendedJSON = userStatisticsJSON + `,"country_rank":11,"rank":{"country":11}`

// decodeStrict fails on any key the target does not map
func decodeStrict(t *testing.T, body string, into any) {
	t.Helper()
	decoder := json.NewDecoder(strings.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(into); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

// assertFilled fails on any field a fully populated response left at its zero value
func assertFilled(t *testing.T, path string, v reflect.Value) {
	t.Helper()
	if v.Type() == reflect.TypeOf(time.Time{}) {
		if v.Interface().(time.Time).IsZero() {
			t.Errorf("%s is zero", path)
		}
		return
	}
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			t.Errorf("%s is nil", path)
			return
		}
		assertFilled(t, path, v.Elem())
	case reflect.Slice:
		if v.Len() == 0 {
			t.Errorf("%s is empty", path)
			return
		}
		assertFilled(t, path+"[0]", v.Index(0))
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			assertFilled(t, path+"."+v.Type().Field(i).Name, v.Field(i))
		}
	default:
		if v.IsZero() {
			t.Errorf("%s is zero", path)
		}
	}
}

func TestUserStatisticsMapsEveryKey(t *testing.T) {
	var stats UserStatisticsExtended
	decodeStrict(t, `{`+userStatisticsExtendedJSON+`}`, &stats)
	assertFilled(t, "statistics", reflect.ValueOf(stats))
}

func TestGetUsersFillsEveryField(t *testing.T) {
	var entry = `{` + userJSON + `,` + userAttributesJSON + `,"statistics_rulesets":{"osu":{` + userStatisticsJSON + `}}}`

	var keys UserWithRulesets
	decodeStrict(t, entry, &keys)

	got, err := clientWith(&roundTripFunc{body: `{"users":[` + entry + `]}`}).GetUsers([]int64{2})
	if err != nil {
		t.Fatal(err)
	}
	assertFilled(t, "user", reflect.ValueOf(*got[0]))
}

func TestGetUsersKeepsMissingStatisticsNil(t *testing.T) {
	got, err := clientWith(&roundTripFunc{body: `{"users":[{"id":2,"statistics_rulesets":null}]}`}).GetUsers([]int64{2})
	if err != nil {
		t.Fatal(err)
	}
	if got[0].StatisticsRulesets != nil {
		t.Errorf("rulesets = %+v, want nil", got[0].StatisticsRulesets)
	}

	got, err = clientWith(&roundTripFunc{body: `{"users":[{"id":2,"statistics_rulesets":{"mania":{"pp":1}}}]}`}).GetUsers([]int64{2})
	if err != nil {
		t.Fatal(err)
	}
	if got[0].StatisticsRulesets.Osu != nil {
		t.Errorf("osu = %+v, want nil", got[0].StatisticsRulesets.Osu)
	}
}

func TestUserKeepsNullableFieldsNull(t *testing.T) {
	var body = `{"id":2,"default_group":null,"last_visit":null,"profile_colour":null,"team":null,
		"cover":{"custom_url":null,"id":null,"url":"u"},"statistics":{"global_rank":null,"country_rank":null}}`
	user, err := clientWith(&roundTripFunc{body: body}).GetUser(2)
	if err != nil {
		t.Fatal(err)
	}
	if user.DefaultGroup != nil || user.LastVisit != nil || user.ProfileColour != nil || user.Team != nil {
		t.Errorf("user = %+v, want nullable fields nil", user.User)
	}
	if user.Cover.CustomURL != nil || user.Cover.ID != nil {
		t.Errorf("cover = %+v, want nullable fields nil", user.Cover)
	}
	if user.Statistics.GlobalRank != nil || user.Statistics.CountryRank != nil {
		t.Errorf("statistics = %+v, want ranks nil", user.Statistics)
	}
}
