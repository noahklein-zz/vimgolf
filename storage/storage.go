package storage

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/noahklein/vimgolf/vim"

	"github.com/jinzhu/gorm"
)

type Challenge struct {
	gorm.Model
	Start  string
	Target string
}

type Solution struct {
	gorm.Model
	Command   string
	Score     int `gorm:"default:0"`
	User      string
	Challenge Challenge
}

var db *gorm.DB

var models = []interface{}{
	Challenge{},
	Solution{},
}

// Init initializes the DB and runs migrations.
func Init(driver, addr string) {
	d, err := gorm.Open(driver, addr)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	db = d

	db.AutoMigrate(models...)
}

// GetDB returns the singleton db instance
func GetDB() *gorm.DB {
	return db
}

// ChallengeByID gets a challenge by id.
func ChallengeByID(id uint) (*Challenge, error) {
	var c *Challenge
	if err := db.First(c, id).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func isValidUsername(u string) bool {
	if len(u) < 3 {
		return false
	}
	const regex = `^\w+$`
	ok, err := regexp.MatchString(regex, u)
	if err != nil {
		panic(err)
	}
	return ok
}

func (c Challenge) Answer(ctx context.Context, cmd, user string) (*Solution, error) {
	if !isValidUsername(user) {
		return nil, fmt.Errorf("invalid username %s", user)
	}

	got, err := vim.AttemptChallenge(ctx, c.Start, cmd)
	if err != nil {
		return nil, err
	}

	score := 0
	if got == c.Target {
		score = vim.Score(got)
	}
	solution := &Solution{
		Command:   cmd,
		Score:     score,
		User:      user,
		Challenge: c,
	}
	db.Create(solution)
	return solution, nil
}
