package helpers

import (
	"github.com/HasanShahjahan/go-guest/api/dtos"
	"github.com/HasanShahjahan/go-guest/api/models"
)

func GuestDtoFromEntity(guestEntity models.Guest) dtos.GuestDto {
	return dtos.GuestDto{Name: guestEntity.Name}
}
