package testdata

func A() error {
	return nil
}

func _() {
	var i int
	if i != 0 {
		panic("i not 0")
	}
}

//go:generate go test -coverprofile=profile.out
