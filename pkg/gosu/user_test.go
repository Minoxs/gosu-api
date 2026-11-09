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

const userExtendedJSON = userJSON + `,` + userAttributesJSON + `,"statistics":{` + userStatisticsExtendedJSON + `},
	"account_history":[{"description":"d","id":1,"length":60,"permanent":true,"timestamp":"2026-01-02T03:04:05Z","type":"note"}],
	"active_tournament_banner":{"id":1,"tournament_id":2,"image":"i","image@2x":"i2"},
	"active_tournament_banners":[{"id":1,"tournament_id":2,"image":"i","image@2x":"i2"}],
	"badges":[{"awarded_at":"2026-01-02T03:04:05Z","description":"d","image":"b.png","image@2x_url":"b2","image_url":"b1","url":"u"}],
	"beatmap_playcounts_count":1,"comments_count":2,"cover_url":"https://c",
	"daily_challenge_user_stats":{"daily_streak_best":1,"daily_streak_current":2,"last_update":"2026-01-02T00:00:00Z",
	"last_weekly_streak":"2026-01-02T00:00:00Z","playcount":3,"top_10p_placements":4,"top_50p_placements":5,"user_id":2,
	"weekly_streak_best":6,"weekly_streak_current":7},
	"discord":"d","favourite_beatmapset_count":3,"follower_count":4,"graveyard_beatmapset_count":5,"guest_beatmapset_count":6,
	"has_supported":true,"interests":"i","join_date":"2007-08-28T03:09:12Z","kudosu":{"available":1,"total":2},"location":"l",
	"loved_beatmapset_count":7,"mapping_follower_count":8,
	"matchmaking_stats":[{"first_placements":1,"is_rating_provisional":true,"plays":2,
	"pool":{"active":true,"id":38,"name":"p","ruleset_id":1,"type":"ranked_play","variant_id":1},"pool_id":38,"rank":3,
	"rank_percent":0.03,"rating":999,"recent_history":[{"created_at":"2026-01-02T03:04:05Z","elo_after":999,"id":1,"result":"win","room_id":2}],
	"total_points":4,"user_id":2}],
	"max_blocks":9,"max_friends":10,"monthly_playcounts":[{"start_date":"2007-10-01","count":1}],"nominated_beatmapset_count":11,
	"occupation":"o","page":{"html":"<div></div>","raw":"[b]"},"pending_beatmapset_count":12,"playmode":"osu","playstyle":["mouse"],
	"post_count":13,"previous_usernames":["old"],"profile_hue":252,"profile_order":["me"],
	"rank_highest":{"rank":1,"updated_at":"2017-07-07T00:00:00Z"},"rank_history":{"mode":"osu","data":[1,2]},
	"ranked_and_approved_beatmapset_count":14,"ranked_beatmapset_count":15,
	"replays_watched_counts":[{"start_date":"2012-09-01","count":1}],"scores_best_count":16,"scores_first_count":17,
	"scores_pinned_count":18,"scores_recent_count":19,"support_level":3,"title":"t","title_url":"tu","twitter":"tw",
	"unranked_beatmapset_count":20,"user_achievements":[{"achieved_at":"2026-01-02T03:04:05Z","achievement_id":299}],"website":"w"`

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

func TestGetUserFillsEveryField(t *testing.T) {
	var body = `{` + userExtendedJSON + `}`

	var keys UserExtended
	decodeStrict(t, body, &keys)

	user, err := clientWith(&roundTripFunc{body: body}).GetUser(2)
	if err != nil {
		t.Fatal(err)
	}
	assertFilled(t, "user", reflect.ValueOf(*user))
}

func TestGetOwnUserFillsEveryField(t *testing.T) {
	var body = `{` + userExtendedJSON + `,"is_restricted":true,"session_verified":true,"statistics_rulesets":{"osu":{` + userStatisticsJSON + `}}}`

	var keys OwnUser
	decodeStrict(t, body, &keys)

	user, err := resourceClientWith(&roundTripFunc{body: body}).GetOwnUser()
	if err != nil {
		t.Fatal(err)
	}
	assertFilled(t, "user", reflect.ValueOf(*user))
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
