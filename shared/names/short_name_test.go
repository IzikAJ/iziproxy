package names

import "testing"

func TestShortNameGeneratorFail(t *testing.T) {
	gen := ShortNameGenerator(func(name string) (isUniq bool) {
		return false
	})
	if _, err := gen.Next(); err == nil {
		t.Errorf("should fail")
	}
	if name, _ := gen.Next(); name != "" {
		t.Errorf("should not return name %v", name)
	}
}

func TestShortNameGeneratorOk(t *testing.T) {
	gen := ShortNameGenerator(func(name string) (isUniq bool) {
		return true
	})
	if _, err := gen.Next(); err != nil {
		t.Errorf("should not fail")
	}
	if name, _ := gen.Next(); len(name) != 2 {
		t.Errorf("should return name with min size")
	}
}

func TestShortNameGeneratorAutoSize(t *testing.T) {
	gen := ShortNameGenerator(func(name string) (isUniq bool) {
		return len(name) >= 4
	})
	if _, err := gen.Next(); err != nil {
		t.Errorf("should not fail")
	}
	if name, _ := gen.Next(); len(name) != 4 {
		t.Errorf("should return name with first passed size %v", name)
	}
}

func TestShortNameGeneratorFill(t *testing.T) {
	var names = make(map[string]string)
	gen := ShortNameGenerator(func(name string) bool {
		_, ok := names[name]
		return !ok
	})
	(gen.(*shortName)).maxSize = 3
	for index := 0; index < 1000; index++ {
		if name, err := gen.Next(); err == nil {
			names[name] = name
		} else {
			t.Errorf("should not fail at least (10k times): %v", index)
			break
		}
	}
}
