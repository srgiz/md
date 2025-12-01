package internal

import (
	"md/internal/domain/cmdbus"
	"net/http"
)

func handle[T any](bus *cmdbus.Bus, r *http.Request, cmd any, res *T) error {
	dto, err := bus.Handle(r.Context(), cmd)

	if dto != nil {
		if data := dto.(*T); data != nil {
			*res = *data
		}
	}

	return err
}

/*
func setResDto[T any](res *T, dto any) {
	if dto != nil {
		if data := dto.(*T); data != nil {
			*res = *data
		}
	}
}*/
