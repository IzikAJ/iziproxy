package names

// Generator - generator interface
type sillyName struct {
	missTimes int
	hitTimes  int
}

func (name *sillyName) Next() string {
	panic("not implemented")
	// return "test"
}

func (name *sillyName) Hit() {
	name.hitTimes++
}

func (name *sillyName) Miss() {
	name.missTimes++
}

// SillyName - silly name geerator implementation
var SillyName = sillyName{
	missTimes: 0,
	hitTimes:  0,
}
