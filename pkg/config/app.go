package config

type Target struct {
	// description of this target for logs
	Name string `mapstructure:"name"`
	// pkg or filename to watch
	Pkg string `mapstructure:"pkg"`
	// command to run when pkg changes
	Command string `mapstructure:"command"`
	// args to pass to the command
	Args []string `mapstructure:"args"`
}

type Blob struct {
	Targets []Target `mapstructure:"target"`
	DelayMS int      `mapstructure:"delay_ms"`
}
