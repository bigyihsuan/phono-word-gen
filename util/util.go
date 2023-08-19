package util

import (
	"math/rand"

	"honnef.co/go/js/dom/v2"
)

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
	dom.GetWindow().Console().Call("log", o...)
}
func LogError(o ...any) {
	dom.GetWindow().Console().Call("error", o...)
}

func AnySlice[T any](ele []T) (o []any) {
	for _, e := range ele {
		o = append(o, e)
	}
	return
}
