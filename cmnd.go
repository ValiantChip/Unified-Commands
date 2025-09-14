package cmnd

import (
	"errors"
	"strings"
)

type commandDescription struct {
	name        string
	description string
}

type Handler struct {
	forced              chan []string
	commands            map[string]Command
	commandDescriptions []commandDescription
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

var ErrAlreadyForced = errors.New("response already forced")

func (h *Handler) ForceResponse(r Runner) error {
	if h.forced == nil {
		return ErrAlreadyForced
	}
	h.forced = make(chan []string)
	cmd := <-h.forced
	h.forced = nil
	return r(cmd)
}

func (h *Handler) AddCommand(name, description string, runner Runner) {
	command := Command{
		Description: description,
		Runner:      runner,
	}

	h.commands[name] = command
	h.commandDescriptions = append(h.commandDescriptions, commandDescription{name: name, description: description})
}

func (h *Handler) Handle(input string) (error, bool) {
	args := strings.Split(input, " ")

	return h.HandleArgs(args)
}

func (h *Handler) HandleArgs(args []string) (error, bool) {
	if len(args) == 0 {
		return nil, false
	}

	if h.forced != nil {
		h.forced <- args
		return nil, true
	}

	if cmd, ok := h.commands[args[0]]; ok {
		return cmd.Runner(args), true
	} else {
		return nil, false
	}

}

func (h *Handler) GetDescription() string {
	var desc strings.Builder
	for _, c := range h.commandDescriptions {
		desc.WriteString(c.name)
		if c.description != "" {
			desc.WriteString(":")
			desc.WriteString("\n    \t")
			desc.WriteString(strings.ReplaceAll(c.description, "\n", "\n    \t"))
		}
		desc.WriteString("\n")
	}

	return desc.String()
}

type Runner func([]string) error

type Command struct {
	Description string
	Runner      Runner
}
