package zork

import (
	"math/rand"
	"time"
)

func PickOne(data RndSelect) string {
	if len(data.Unselected) == 0 {
		data.Unselected = data.Selected
		data.Selected = []string{}
	}
	if data.Selected == nil {
		data.Selected = []string{}
	}
	rand.Seed(time.Now().Unix())
	rnd := rand.Intn(len(data.Unselected))
	msg := data.Unselected[rnd]
	data.Selected = append(data.Selected, msg)
	data.Unselected[rnd] = data.Unselected[len(data.Unselected)-1]
	data.Unselected[len(data.Unselected)-1] = ""
	data.Unselected = data.Unselected[:len(data.Unselected)-1]
	return msg
}

func Prob(base int, isLooser bool) bool {
	if isLooser {
		return Zprob(base)
	}
	rand.Seed(time.Now().Unix())
	return base > rand.Intn(100)
}

func Zprob(base int) bool {
	rand.Seed(time.Now().Unix())
	if Lucky {
		return base > rand.Intn(100)
	}
	return base > rand.Intn(300)
}
