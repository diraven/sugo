package publicroles

import (
	"github.com/diraven/sugo"
)

// Init initializes module.
func Init(sg *sugo.Instance) {
	sg.AddCommand(cmd)

	sg.AddStartupHandler(func(sg *sugo.Instance) error {
		err := storage.load()
		if err != nil {
			return err
		}
		return nil
	})

	sg.AddShutdownHandler(func(sg *sugo.Instance) error {
		err := storage.save()
		if err != nil {
			return err
		}
		return nil
	})
}
