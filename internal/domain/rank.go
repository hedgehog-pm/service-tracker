package domain

import "math"

const rankStep = 1024

func RankBetween(previous, next *int) (int, error) {
	switch {
	case previous == nil && next == nil:
		return rankStep, nil
	case previous == nil:
		if *next < math.MinInt+rankStep {
			return 0, ErrRankExhausted
		}

		return *next - rankStep, nil
	case next == nil:
		if *previous > math.MaxInt-rankStep {
			return 0, ErrRankExhausted
		}

		return *previous + rankStep, nil
	}

	if *previous >= *next {
		return 0, ErrRankExhausted
	}

	var rank int
	if *previous < 0 && *next >= 0 {
		rank = (*previous + *next) / 2
	} else {
		rank = *previous + (*next-*previous)/2
	}

	if rank <= *previous || rank >= *next {
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
