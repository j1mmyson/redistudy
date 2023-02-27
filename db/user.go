package db

import (
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

type User struct {
	UserID string `json:"uid" binding:"required"`
	Score  int    `json:"score" binding:"required"`
	Rank   int    `json:"rank"`
}

func (db *Database) SaveUser(user *User) error {
	member := redis.Z{
		Score:  float64(user.Score),
		Member: user.UserID,
	}
	pipe := db.Client.TxPipeline()
	pipe.ZAdd(Ctx, leaderboardKey, member)
	rank := pipe.ZRank(Ctx, leaderboardKey, user.UserID)
	_, err := pipe.Exec(Ctx)
	if err != nil {
		return err
	}
	fmt.Println(rank.Val(), err)
	user.Rank = int(rank.Val())
	return nil
}

func (db *Database) GetUser(userID string) (*User, error) {
	pipe := db.Client.TxPipeline()
	score := pipe.ZScore(Ctx, leaderboardKey, userID)
	rank := pipe.ZRank(Ctx, leaderboardKey, userID)
	_, err := pipe.Exec(Ctx)
	if err != nil {
		return nil, err
	}

	if score == nil {
		return nil, ErrNil
	}

	return &User{
		UserID: userID,
		Score:  int(score.Val()),
		Rank:   int(rank.Val()),
	}, nil
}
