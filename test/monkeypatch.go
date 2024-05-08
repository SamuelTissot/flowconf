package test

func MonkeyPatch[T any](this *T, toThat *T) func() {
	original := *this
	*this = *toThat

	return func() {
		*this = original
	}
}
