package testdata

func A() error {
	return nil
}

//go:generate go test -coverprofile=profile.out
