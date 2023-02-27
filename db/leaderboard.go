package db

var leaderboardKey = "leaderboard"

type Leaderboard struct {
	Count int `json:"count"`
	Users []*User
}

func (db *Database) GetLeaderboard() (*Leaderboard, error) {
	scores := db.Client.ZRevRangeWithScores(Ctx, leaderboardKey, 0, -1)
	if scores == nil {
		return nil, ErrNil
	}

	count := len(scores.Val())
	users := make([]*User, count)
	for idx, member := range scores.Val() {
		users[idx] = &User{
			UserID: member.Member.(string),
			Score:  int(member.Score),
			Rank:   idx + 1,
		}
	}
	leaderboard := &Leaderboard{
		Count: count,
		Users: users,
	}

	return leaderboard, nil
}
