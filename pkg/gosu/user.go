package gosu

import "time"

// Ranks is osu!'s nested rank block
type Ranks struct {
	Country *int `json:"country"`
}

// Level is a user's level and the percentage progress towards the next one
type Level struct {
	Current  int     `json:"current"`
	Progress float64 `json:"progress"`
}

// GradeCounts is how many ranked scores a user holds at each grade
type GradeCounts struct {
	A   int `json:"a"`
	S   int `json:"s"`
	SH  int `json:"sh"`
	SS  int `json:"ss"`
	SSH int `json:"ssh"`
}

// UserStatistics is a user's osu!standard ranking summary as every user endpoint
// reports it
type UserStatistics struct {
	Accuracy float64 `json:"accuracy"`
	// Deprecated: use Accuracy, which is a fraction in [0, 1] where this is a percentage
	HitAccuracy            float64     `json:"hit_accuracy"`
	Count100               int         `json:"count_100"`
	Count300               int         `json:"count_300"`
	Count50                int         `json:"count_50"`
	CountMiss              int         `json:"count_miss"`
	GlobalRank             *int        `json:"global_rank"`
	GlobalRankPercent      *float64    `json:"global_rank_percent"`
	GradeCounts            GradeCounts `json:"grade_counts"`
	IsRanked               bool        `json:"is_ranked"`
	Level                  Level       `json:"level"`
	MaximumCombo           int         `json:"maximum_combo"`
	PlayCount              int         `json:"play_count"`
	PlayTime               int         `json:"play_time"`
	PP                     float64     `json:"pp"`
	RankedScore            int64       `json:"ranked_score"`
	ReplaysWatchedByOthers int         `json:"replays_watched_by_others"`
	TotalHits              int64       `json:"total_hits"`
	TotalScore             int64       `json:"total_score"`
}

// UserStatisticsExtended is UserStatistics plus the country standing the single-user
// endpoints add
type UserStatisticsExtended struct {
	UserStatistics
	CountryRank *int  `json:"country_rank"`
	Rank        Ranks `json:"rank"`
}

// Country is an osu! user's country
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Cover is an osu! user's profile cover
type Cover struct {
	CustomURL *string `json:"custom_url"`
	ID        *string `json:"id"`
	URL       string  `json:"url"`
}

// UserGroup is a user's membership of an osu! group
type UserGroup struct {
	Colour         *string  `json:"colour"`
	HasListing     bool     `json:"has_listing"`
	HasPlaymodes   bool     `json:"has_playmodes"`
	ID             int      `json:"id"`
	Identifier     string   `json:"identifier"`
	IsProbationary bool     `json:"is_probationary"`
	Name           string   `json:"name"`
	Playmodes      []string `json:"playmodes"`
	ShortName      string   `json:"short_name"`
}

// Team is the osu! team a user belongs to
type Team struct {
	FlagURL   *string `json:"flag_url"`
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	ShortName string  `json:"short_name"`
}

// User is the profile every user endpoint returns
type User struct {
	ID            int64       `json:"id"`
	Username      string      `json:"username"`
	AvatarURL     string      `json:"avatar_url"`
	CountryCode   string      `json:"country_code"`
	DefaultGroup  *string     `json:"default_group"`
	IsActive      bool        `json:"is_active"`
	IsBot         bool        `json:"is_bot"`
	IsDeleted     bool        `json:"is_deleted"`
	IsOnline      bool        `json:"is_online"`
	IsSupporter   bool        `json:"is_supporter"`
	LastVisit     *time.Time  `json:"last_visit"`
	PMFriendsOnly bool        `json:"pm_friends_only"`
	ProfileColour *string     `json:"profile_colour"`
	Country       Country     `json:"country"`
	Cover         Cover       `json:"cover"`
	Groups        []UserGroup `json:"groups"`
	Team          *Team       `json:"team"`
}

// UserStatisticsRulesets is a user's statistics by ruleset. Osu is nil for an account
// osu! holds no osu!standard statistics for
type UserStatisticsRulesets struct {
	Osu *UserStatistics `json:"osu"`
}

// UserWithRulesets is the profile the bulk users endpoint returns. StatisticsRulesets is
// nil for an account osu! holds no statistics for
type UserWithRulesets struct {
	User
	StatisticsRulesets *UserStatisticsRulesets `json:"statistics_rulesets"`
}

// UserExtended is the full osu! profile the single-user endpoints return
type UserExtended struct {
	User
	// Deprecated: use Cover.URL
	CoverURL   string                 `json:"cover_url"`
	Statistics UserStatisticsExtended `json:"statistics"`
	JoinDate   time.Time              `json:"join_date"`
}
