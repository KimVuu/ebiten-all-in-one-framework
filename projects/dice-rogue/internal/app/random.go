package app

import "math/rand"

type RandomSource struct {
	rng      *rand.Rand
	scripted []int
}

func NewRandomSource(seed int64) *RandomSource {
	return &RandomSource{
		rng: rand.New(rand.NewSource(seed)),
	}
}

func newRandomSourceWithScript(seed int64, values ...int) *RandomSource {
	source := NewRandomSource(seed)
	source.scripted = append(source.scripted, values...)
	return source
}

func (source *RandomSource) NextInt(n int) int {
	if n <= 0 {
		return 0
	}
	if source == nil {
		return 0
	}
	if len(source.scripted) > 0 {
		value := source.scripted[0]
		source.scripted = source.scripted[1:]
		if value < 0 {
			value = -value
		}
		return value % n
	}
	return source.rng.Intn(n)
}
