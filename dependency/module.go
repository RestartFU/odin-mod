package dependency

import "github.com/RestartFU/gophig"

type Module struct {
	Version      string   `toml:"version"`
	Dependencies []string `toml:"dependencies"`
}

func ReadModule(path string) (Module, error) {
	mod := Module{}
	err := gophig.GetConfComplex(path+"/odin.toml", gophig.TOMLMarshaler{}, &mod)
	return mod, err
}

func (m Module) Save() error {
	err := gophig.SetConfComplex("odin.toml", gophig.TOMLMarshaler{}, m, 0664)
	return err
}
