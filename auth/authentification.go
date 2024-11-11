package authentification

import (
	"context"
	"log"
	"slices"
)

type Privilege string

type privilegeKey struct{}
type checkedKey struct{}

func Check(ctx context.Context, ps ...Privilege) (_ context.Context, ok bool) {
	granted, ok := ctx.Value(privilegeKey{}).([]Privilege)
	if !ok {
		log.Printf("Privilege key not found")
		return ctx, false
	}

	for _, p := range ps {
		if !slices.Contains(granted, p) {
			log.Printf("Granted: %v, Requires: %v", granted, p)
			return ctx, false
		}
	}

	return context.WithValue(ctx, checkedKey{}, true), true
}

func Grant(ctx context.Context, ps ...Privilege) context.Context {
	if ctx.Value(privilegeKey{}) != nil {
		panic("Grant called multiple times!")
	}

	return context.WithValue(ctx, privilegeKey{}, ps)
}

func Must(ctx context.Context) (ok bool) {
	if ctx.Value(checkedKey{}) == nil {
		return false
	}

	return true
}
