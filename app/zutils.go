package app

import "log"

// I just wanted to check how many would actually look into the `must` function.
// I'm just lazy and this is a dummy project, let me write some horrible code now
// that I can!

func must[T any](t T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return t
}
