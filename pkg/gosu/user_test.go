package gosu

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

// The fixture repeats the keys of a live response, every value non-zero so assertFilled
// can tell a mapped field from an unmapped one
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
	var stats UserStatistics
	decodeStrict(t, `{`+userStatisticsExtendedJSON+`}`, &stats)
	assertFilled(t, "statistics", reflect.ValueOf(stats))
}
