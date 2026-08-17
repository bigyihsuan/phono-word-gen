package util

import (
	"encoding/json"
	"math/rand/v2"

	"honnef.co/go/js/dom/v2"
)

// a fixed random source for unit tests
func RandomSource() *rand.Rand {
	return rand.New(rand.NewPCG(0, 0))
}

func PeakedPowerLaw(max, mode, prob int, rs *rand.Rand) int {
	if RandomPercentage(rs) < 50 {
		return mode + PowerLaw(max-mode, prob, rs)
	}
	return mode + PowerLaw(mode+1, prob, rs)
}

func PowerLaw(max, percentage int, rs *rand.Rand) int {
	for r := 0; ; r = (r + 1) % max {
		if RandomPercentage(rs) < percentage {
			return r
		}
	}
}

func RandomPercentage(rs *rand.Rand) int {
	return rs.IntN(101) + 1
}

func Log(o ...any) {
	args := []any{}
	for _, arg := range o {
		m, err := ToMap(arg)
		if err != nil {
			LogError(err.Error())
			continue
		}
		args = append(args, m)
	}
	dom.GetWindow().Console().Call("log", args...)
}
func LogError(o ...any) {
	args := []any{}
	for _, arg := range o {
		m, err := ToMap(arg)
		if err != nil {
			LogError(err.Error())
			continue
		}
		args = append(args, m)
	}
	dom.GetWindow().Console().Call("error", args...)
}

func ToMap(x any) (map[string]any, error) {
	data, err := json.Marshal(struct{ V any }{x})
	if err != nil {
		return nil, err
	}
	m := map[string]any{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
