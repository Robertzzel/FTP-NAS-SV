package commands

type Command interface {
	Execute() (int, error)
}
