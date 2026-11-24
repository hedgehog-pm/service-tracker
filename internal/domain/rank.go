package domain

import (
	"math"

	"github.com/samber/lo"
)

const rankStep = 1024

func RankBetween(previous, next *int) (int, error) {
	previousRank := lo.FromPtr(previous)
	nextRank := lo.FromPtr(next)

	switch {
	case previous == nil && next == nil:
		return rankStep, nil
	case previous == nil:
		if nextRank < math.MinInt+rankStep {
			return 0, ErrRankExhausted
		}

		return nextRank - rankStep, nil
	case next == nil:
		if previousRank > math.MaxInt-rankStep {
			return 0, ErrRankExhausted
		}

		return previousRank + rankStep, nil
	}

	if previousRank >= nextRank {
		return 0, ErrRankExhausted
	}

	var rank int
	if previousRank < 0 && nextRank >= 0 {
		rank = (previousRank + nextRank) / 2
	} else {
		rank = previousRank + (nextRank-previousRank)/2
	}

	if rank <= previousRank || rank >= nextRank {
		return 0, ErrRankExhausted
	}

	return rank, nil
}

func RebalanceTaskRanks(tasks []Task) error {
	if len(tasks) > math.MaxInt/rankStep {
		return ErrRankExhausted
	}

	for i := range tasks {
		tasks[i].Rank = (i + 1) * rankStep
	}

	return nil
}
