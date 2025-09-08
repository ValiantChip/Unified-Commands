package cmnd

import "strings"

type Handler struct {
	commands map[string]Command
}

type HandlerArg struct {
	Name        string
	Description string
	Runner      Runner
}

func NewHandler(commands ...HandlerArg) *Handler {
	handler := &Handler{
		commands: make(map[string]Command),
	}

	for _, command := range commands {
		handler.AddCommand(command.Name, command.Description, command.Runner)
	}

	return handler
}

func (h *Handler) AddCommand(name, description string, runner Runner) {
	command := Command{
		Description: description,
		Runner:      runner,
	}

	h.commands[name] = command
}

func (h *Handler) Handle(input string) (error, bool) {
	args := strings.Split(input, " ")

	return h.HandleArgs(args)
}

func (h *Handler) HandleArgs(args []string) (error, bool) {
	if len(args) == 0 {
		return nil, false
	}

	if cmd, ok := h.commands[args[0]]; ok {
		return cmd.Runner(args), true
	} else {
		return nil, false
	}

}

func (h *Handler) GetDescription() string {
	var desc strings.Builder
	for name, cmd := range h.commands {
		desc.WriteString(name)
		desc.WriteString("\n    \t")
		desc.WriteString(strings.ReplaceAll(cmd.Description, "\n", "\n    \t"))
	}

	return desc.String()
}

type Runner func([]string) error

type Command struct {
	Description string
	Runner      Runner
}
