package device

import (
	"fmt"
	"time"
)

type Device struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *Device) String() string {
	return fmt.Sprintf("Device<%s %s>", d.ID, d.Name)
}

type CreateDeviceRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
