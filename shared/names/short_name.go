package names

import "math/rand"

// Generator - generator interface
type shortName struct {
	minSize   int
	maxSize   int
	missTimes int
	checker   func(string) bool
}

const missRate = 10

var symbols = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func (name *shortName) Next() (ans string, err error) {
	for {
		ans = name.make()
		if name.checker(ans) {
			name.hit()
			break
		} else {
			err = name.miss()
			if err != nil {
				ans = ""
				break
			}
		}
	}
	return
}

func (name *shortName) make() string {
	mass := make([]rune, name.minSize)
	for i := range mass {
		mass[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(mass)
}

func (name *shortName) hit() {
	name.missTimes = 0
}

func (name *shortName) miss() (err error) {
	name.missTimes++
	if name.missTimes > (name.minSize * missRate) {
		if name.minSize >= name.maxSize {
			return MissError
		}
		name.minSize++
		name.missTimes = 0
	}
	return
}

// ShortNameGenerator - create short name genarator
func ShortNameGenerator(checker func(string) bool) Generator {
	return &shortName{
		minSize: 2,
		maxSize: 5,
		checker: checker,
	}
}
