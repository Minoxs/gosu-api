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

// UserAccountHistory is one moderation event on a user's account. Length is in seconds
type UserAccountHistory struct {
	Description *string   `json:"description"`
	ID          int64     `json:"id"`
	Length      int       `json:"length"`
	Permanent   bool      `json:"permanent"`
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`
}

// ProfileBanner is a tournament banner shown on a user's profile
type ProfileBanner struct {
	ID           int64   `json:"id"`
	TournamentID int64   `json:"tournament_id"`
	Image        *string `json:"image"`
	Image2x      *string `json:"image@2x"`
}

// UserBadge is a badge awarded to a user
type UserBadge struct {
	AwardedAt   time.Time `json:"awarded_at"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Image2xURL  string    `json:"image@2x_url"`
	ImageURL    string    `json:"image_url"`
	URL         string    `json:"url"`
}

// Kudosu is a user's modding currency
type Kudosu struct {
	Available int `json:"available"`
	Total     int `json:"total"`
}

// UserPage is a user's profile page body
type UserPage struct {
	HTML string `json:"html"`
	Raw  string `json:"raw"`
}

// UserMonthlyPlaycount is one month of a per-month counter. StartDate is YYYY-MM-DD,
// which time.Time cannot decode
type UserMonthlyPlaycount struct {
	StartDate string `json:"start_date"`
	Count     int    `json:"count"`
}

// RankHighest is the best global rank a user has reached
type RankHighest struct {
	Rank      int       `json:"rank"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RankHistory is a user's recent daily global ranks, oldest first
type RankHistory struct {
	Mode string `json:"mode"`
	Data []int  `json:"data"`
}

// UserAchievement is a medal a user has unlocked
type UserAchievement struct {
	AchievedAt    time.Time `json:"achieved_at"`
	AchievementID int       `json:"achievement_id"`
}

// DailyChallengeUserStats is a user's daily challenge record
type DailyChallengeUserStats struct {
	DailyStreakBest     int       `json:"daily_streak_best"`
	DailyStreakCurrent  int       `json:"daily_streak_current"`
	LastUpdate          time.Time `json:"last_update"`
	LastWeeklyStreak    time.Time `json:"last_weekly_streak"`
	Playcount           int       `json:"playcount"`
	Top10pPlacements    int       `json:"top_10p_placements"`
	Top50pPlacements    int       `json:"top_50p_placements"`
	UserID              int64     `json:"user_id"`
	WeeklyStreakBest    int       `json:"weekly_streak_best"`
	WeeklyStreakCurrent int       `json:"weekly_streak_current"`
}

// MatchmakingPool is a matchmaking pool
type MatchmakingPool struct {
	Active    bool   `json:"active"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	RulesetID int    `json:"ruleset_id"`
	Type      string `json:"type"`
	VariantID int    `json:"variant_id"`
}

// MatchmakingResult is one recent matchmaking room a user played
type MatchmakingResult struct {
	CreatedAt time.Time `json:"created_at"`
	EloAfter  int       `json:"elo_after"`
	ID        int64     `json:"id"`
	Result    string    `json:"result"`
	RoomID    int64     `json:"room_id"`
}

// MatchmakingStats is a user's standing in one matchmaking pool
type MatchmakingStats struct {
	FirstPlacements     int                 `json:"first_placements"`
	IsRatingProvisional bool                `json:"is_rating_provisional"`
	Plays               int                 `json:"plays"`
	Pool                MatchmakingPool     `json:"pool"`
	PoolID              int64               `json:"pool_id"`
	Rank                int                 `json:"rank"`
	RankPercent         float64             `json:"rank_percent"`
	Rating              int                 `json:"rating"`
	RecentHistory       []MatchmakingResult `json:"recent_history"`
	TotalPoints         int                 `json:"total_points"`
	UserID              int64               `json:"user_id"`
}

// UserExtended is the full osu! profile the single-user endpoints return
type UserExtended struct {
	User
	AccountHistory []UserAccountHistory `json:"account_history"`
	// Deprecated: use ActiveTournamentBanners
	ActiveTournamentBanner  *ProfileBanner  `json:"active_tournament_banner"`
	ActiveTournamentBanners []ProfileBanner `json:"active_tournament_banners"`
	Badges                  []UserBadge     `json:"badges"`
	BeatmapPlaycountsCount  int             `json:"beatmap_playcounts_count"`
	CommentsCount           int             `json:"comments_count"`
	// Deprecated: use Cover.URL
	CoverURL                         string                  `json:"cover_url"`
	DailyChallengeUserStats          DailyChallengeUserStats `json:"daily_challenge_user_stats"`
	Discord                          *string                 `json:"discord"`
	FavouriteBeatmapsetCount         int                     `json:"favourite_beatmapset_count"`
	FollowerCount                    int                     `json:"follower_count"`
	GraveyardBeatmapsetCount         int                     `json:"graveyard_beatmapset_count"`
	GuestBeatmapsetCount             int                     `json:"guest_beatmapset_count"`
	HasSupported                     bool                    `json:"has_supported"`
	Interests                        *string                 `json:"interests"`
	JoinDate                         time.Time               `json:"join_date"`
	Kudosu                           Kudosu                  `json:"kudosu"`
	Location                         *string                 `json:"location"`
	LovedBeatmapsetCount             int                     `json:"loved_beatmapset_count"`
	MappingFollowerCount             int                     `json:"mapping_follower_count"`
	MatchmakingStats                 []MatchmakingStats      `json:"matchmaking_stats"`
	MaxBlocks                        int                     `json:"max_blocks"`
	MaxFriends                       int                     `json:"max_friends"`
	MonthlyPlaycounts                []UserMonthlyPlaycount  `json:"monthly_playcounts"`
	NominatedBeatmapsetCount         int                     `json:"nominated_beatmapset_count"`
	Occupation                       *string                 `json:"occupation"`
	Page                             UserPage                `json:"page"`
	PendingBeatmapsetCount           int                     `json:"pending_beatmapset_count"`
	Playmode                         string                  `json:"playmode"`
	Playstyle                        []string                `json:"playstyle"`
	PostCount                        int                     `json:"post_count"`
	PreviousUsernames                []string                `json:"previous_usernames"`
	ProfileHue                       *int                    `json:"profile_hue"`
	ProfileOrder                     []string                `json:"profile_order"`
	RankHighest                      *RankHighest            `json:"rank_highest"`
	RankHistory                      *RankHistory            `json:"rank_history"`
	RankedAndApprovedBeatmapsetCount int                     `json:"ranked_and_approved_beatmapset_count"`
	RankedBeatmapsetCount            int                     `json:"ranked_beatmapset_count"`
	ReplaysWatchedCounts             []UserMonthlyPlaycount  `json:"replays_watched_counts"`
	ScoresBestCount                  int                     `json:"scores_best_count"`
	ScoresFirstCount                 int                     `json:"scores_first_count"`
	ScoresPinnedCount                int                     `json:"scores_pinned_count"`
	ScoresRecentCount                int                     `json:"scores_recent_count"`
	Statistics                       UserStatisticsExtended  `json:"statistics"`
	SupportLevel                     int                     `json:"support_level"`
	Title                            *string                 `json:"title"`
	TitleURL                         *string                 `json:"title_url"`
	Twitter                          *string                 `json:"twitter"`
	UnrankedBeatmapsetCount          int                     `json:"unranked_beatmapset_count"`
	UserAchievements                 []UserAchievement       `json:"user_achievements"`
	Website                          *string                 `json:"website"`
}

// OwnUser is the profile of the user a resource token belongs to
type OwnUser struct {
	UserExtended
	IsRestricted       *bool                   `json:"is_restricted"`
	SessionVerified    bool                    `json:"session_verified"`
	StatisticsRulesets *UserStatisticsRulesets `json:"statistics_rulesets"`
}
