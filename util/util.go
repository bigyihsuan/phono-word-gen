package util

import (
	"encoding/json"
	"math/rand"

	"honnef.co/go/js/dom/v2"
)

func PeakedPowerLaw(max, mode, prob int) int {
	if RandomPercentage() < 50 {
		return mode + PowerLaw(max-mode, prob)
	}
	return mode + PowerLaw(mode+1, prob)
}

func PowerLaw(max, percentage int) int {
	for r := 0; ; r = (r + 1) % max {
		if RandomPercentage() < percentage {
			return r
		}
	}
}

func RandomPercentage() int {
	return rand.Intn(101) + 1
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
