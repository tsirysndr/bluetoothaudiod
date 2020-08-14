package utils

import (
	"github.com/godbus/dbus"
	"github.com/pkg/errors"

	bluez "github.com/tsirysndr/bluetoothaudiod/server/bluetooth"
)

func NewBluez() (*bluez.Bluez, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create dbus system bus:")
	}
	b := bluez.NewBluez(conn)
	if err := b.PopulateCache(); err != nil {
		return nil, errors.Wrapf(err, "unable to populate cache")
	}
	return b, nil
}
