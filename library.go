package main

import "strings"

type Library interface {
	Commands() map[string]*Command
	AddCommand(newCmd *Command)
	RemoveCommand(toRemove *Command) bool
	FindCommand(search string) *Command
	Combat() string
}

type MudLib struct {
	commands map[string]*Command
	combat   string
	world    *World
}

func NewLibrary(world *World) *MudLib {
	commands := make(map[string]*Command)
	var combat string
	return &MudLib{
		commands: commands,
		combat:   combat,
		world:    world,
	}
}

func (l *MudLib) Commands() map[string]*Command {
	return l.commands
}

func (l *MudLib) AddCommand(newCmd *Command) {
	_, exists := l.commands[newCmd.Trigger()]
	if !exists {
		l.commands[newCmd.Trigger()] = newCmd
	}
}

func (l *MudLib) RemoveCommand(toRemove *Command) bool {
	_, exists := l.commands[toRemove.Trigger()]
	delete(l.commands, toRemove.Trigger())
	return exists
}

func (l *MudLib) FindCommand(search string) *Command {
	var cmd *Command
	search = strings.ToLower(search)

	// find exact match
	foundCmd, exists := l.commands[search]
	if exists {
		return foundCmd
	}

	// find partial match
	for trigger, command := range l.Commands() {
		if strings.HasPrefix(strings.ToLower(trigger), search) {
			return command
		}
	}

	return cmd
}

func (l *MudLib) LoadCommands() {
	SayCmd := Command{
		trigger:    "say",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&SayCmd)

	LookCmd := Command{
		trigger:    "look",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&LookCmd)

	RenameCmd := Command{
		trigger:    "rename",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&RenameCmd)

	RerollCmd := Command{
		trigger:    "reroll",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&RerollCmd)

	SpawnCmd := Command{
		trigger:    "spawn",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&SpawnCmd)

	StatsCmd := Command{
		trigger:    "stats",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&StatsCmd)

	HealthCmd := Command{
		trigger:    "health",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&HealthCmd)

	ReleaseCmd := Command{
		trigger:    "release",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&ReleaseCmd)

	AttackCmd := Command{
		trigger:    "hit",
		timer:      2000,
		timerType:  TIMER_ATTACK,
		costType:   COST_FAT,
		useCost:    3,
		checkTimer: true,
	}
	l.AddCommand(&AttackCmd)

	InfoCmd := Command{
		trigger:    "information",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&InfoCmd)

	QuitCmd := Command{
		trigger:    "quit",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&QuitCmd)

	InventoryCmd := Command{
		trigger:    "inventory",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&InventoryCmd)

	WealthCmd := Command{
		trigger:    "wealth",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&WealthCmd)

	CreateCmd := Command{
		trigger:    "create",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&CreateCmd)

	GetCmd := Command{
		trigger:    "get",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&GetCmd)

	DropCmd := Command{
		trigger:    "drop",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: true,
	}
	l.AddCommand(&DropCmd)

	GiveCmd := Command{
		trigger:    "give",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&GiveCmd)

	GoCmd := Command{
		trigger:    "go",
		timer:      0,
		timerType:  TIMER_NONE,
		costType:   COST_NONE,
		useCost:    0,
		checkTimer: true,
	}
	l.AddCommand(&GoCmd)
}

func (l *MudLib) Combat() string {
	return l.combat
}
