package cli

type CliCommand interface {
	Exec()
	Help()
}
