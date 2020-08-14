package server

import (
	context "context"
	"fmt"
	"os"

	bluez "github.com/tsirysndr/bluetoothaudiod/server/bluetooth"

	"github.com/godbus/dbus"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/ulule/deepcopier"
)

type Server struct {
	B       *bluez.Bluez
	Adapter string
}

func (s *Server) ListDevices(ctx context.Context, req *Empty) (res *Devices, err error) {
	devices := []*Device{}
	s.B.PopulateCache()
	for _, item := range s.B.Devices {
		var device Device
		deepcopier.Copy(&item).To(&device)
		device.Blocked = item.Blocked
		device.Connected = item.Connected
		devices = append(devices, &device)
	}
	return &Devices{Devices: devices}, nil
}

func (s *Server) ListAdapters(ctx context.Context, req *Empty) (res *Adapters, err error) {
	adapters := []*Adapter{}
	s.B.PopulateCache()
	for _, item := range s.B.Adapters {
		var adapter Adapter
		deepcopier.Copy(&item).To(&adapter)
		adapters = append(adapters, &adapter)
	}
	return &Adapters{Adapters: adapters}, nil
}

func (s *Server) Connect(ctx context.Context, req *Params) (res *Adapters, err error) {
	adapters := []*Adapter{}
	if err := s.B.Connect(s.Adapter, req.Name); err != nil {
		return &Adapters{Adapters: adapters}, err
	}
	for _, item := range s.B.Adapters {
		var adapter Adapter
		deepcopier.Copy(&item).To(&adapter)
		adapters = append(adapters, &adapter)
	}
	return &Adapters{Adapters: adapters}, nil
}

func (s *Server) Disconnect(ctx context.Context, req *Params) (res *Adapters, err error) {
	adapters := []*Adapter{}
	if err := s.B.Disconnect(s.Adapter, req.Name); err != nil {
		return &Adapters{Adapters: adapters}, err
	}
	s.B.PopulateCache()
	for _, item := range s.B.Adapters {
		var adapter Adapter
		deepcopier.Copy(&item).To(&adapter)
		adapter.Discoverable = item.Discoverable
		adapters = append(adapters, &adapter)
	}
	return &Adapters{Adapters: adapters}, nil
}

func (s *Server) Pair(ctx context.Context, req *Params) (res *Adapters, err error) {
	adapters := []*Adapter{}
	if err := s.B.Pair(s.Adapter, req.Name); err != nil {
		return &Adapters{Adapters: adapters}, err
	}
	s.B.PopulateCache()
	for _, item := range s.B.Adapters {
		var adapter Adapter
		deepcopier.Copy(&item).To(&adapter)
		adapter.Discoverable = item.Discoverable
		adapters = append(adapters, &adapter)
	}
	return &Adapters{Adapters: adapters}, nil
}

func (s *Server) EnableCard(ctx context.Context, req *Card) (res *Status, err error) {
	if req == nil {
		return &Status{Ok: false}, nil
	}

	const (
		PREFIX   = "pcm.!default {\n\ttype asym\n\tplayback.pcm {\n\t\ttype plug\n\t\tslave.pcm \"output\"\n\t}\n\tcapture.pcm {\n\t\ttype plug\n\t\tslave.pcm \"input\"\n\t}\n}"
		OUTPUT_A = "\npcm.output {\n\ttype hw\n\tcard %d\n}"
		INPUT_A  = "\npcm.input {\n\ttype hw\n\tcard %d\n}"
		CTL_A    = "\nctl.!default {\n\ttype hw\n\tcard %d\n}"
		OUTPUT_B = "\npcm.output {\n\ttype bluealsa\n\tdevice \"%s\"\n\tprofile \"a2dp\"\n}"
		INPUT_B  = "\npcm.input {\n\ttype bluealsa\n\tdevice \"%s\"\n\tprofile \"sco\"\n}"
		CTL_B    = "\nctl.!default {\n\ttype bluealsa\n}"
	)

	home, _ := homedir.Dir()
	f, err := os.Create(home + "/.asoundrc")

	defer f.Close()

	if err != nil {
		return &Status{Ok: false}, err
	}

	output := fmt.Sprintf(OUTPUT_A+CTL_A, req.Num, req.Num)

	if req.Address != "" {
		output = fmt.Sprintf(OUTPUT_B+CTL_B, req.Address)
	}

	config := fmt.Sprintf("%s %s", PREFIX, output)

	fmt.Println(config)

	if _, err := f.WriteString(config); err != nil {
		f.Sync()
		return &Status{Ok: false}, err
	}

	return &Status{Ok: true}, nil
}

func (s Server) StartDiscovery(ctx context.Context, req *Adapter) (res *Status, err error) {
	adapter := "hci0"
	if req.Name != "" {
		adapter = req.Name
	}
	go func() {
		if err := s.B.StartDiscovery(adapter); err != nil {
			fmt.Println(err)
			return
		}
		signalChan := s.B.WatchSignal()
		for signal := range signalChan {
			fmt.Printf("received signal=%s => (%d)%v\n", signal.Name, len(signal.Body), signal.Body)
			if signal.Name == "org.freedesktop.DBus.ObjectManager.InterfacesAdded" {
				if len(signal.Body) != 2 {
					continue
				}
				devicePath, ok := signal.Body[0].(dbus.ObjectPath)
				if !ok {
					fmt.Printf("unable to cast '%#v' to dbus.ObjectPath", signal.Body[0])
					continue
				}
				deviceMap, ok := signal.Body[1].(map[string]map[string]dbus.Variant)
				if !ok {
					fmt.Printf("unable to cast '%#v' to device map[string]dbus.Variant", signal.Body[1])
					continue
				}
				devices := s.B.ConvertToDevices(string(devicePath), deviceMap)
				for _, d := range devices {
					fmt.Printf("name=%q alias=%q address=%q, adapter=%q paired=%t connected=%t trusted=%t blocked=%t\n", d.Name, d.Alias, d.Address, d.Adapter, d.Paired, d.Connected, d.Trusted, d.Blocked)
				}
			}
		}
	}()
	return &Status{Ok: true}, nil
}

func (s Server) StopDiscovery(ctx context.Context, req *Adapter) (res *Status, err error) {
	adapter := s.Adapter
	if req.Name != "" {
		adapter = req.Name
	}
	if err := s.B.StopDiscovery(adapter); err != nil {
		fmt.Println(err)
		return &Status{Ok: false}, err
	}
	return &Status{Ok: true}, nil
}
