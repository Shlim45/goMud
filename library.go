package main

import "strings"

type Library interface {
	Commands() map[string]*Command
	AddCommand(newCmd *Command)
	RemoveCommand(toRemove *Command) bool
	FindCommand(search string) *Command

	CharClasses() map[string]CharClass
	AddCharClass(newClass CharClass)
	RemoveCharClass(toRemove CharClass) bool
	FindCharClass(search string) CharClass

	Races() map[string]Race
	AddRace(newClass Race)
	RemoveRace(toRemove Race) bool
	FindRace(search string) Race

	Combat() string
}

type MudLib struct {
	commands    map[string]*Command
	charClasses map[string]CharClass
	races       map[string]Race
	combat      string
	world       *World
}

func NewLibrary(world *World) *MudLib {
	var combat string
	return &MudLib{
		commands:    make(map[string]*Command),
		charClasses: make(map[string]CharClass),
		races:       make(map[string]Race),
		combat:      combat,
		world:       world,
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

func (l *MudLib) CharClasses() map[string]CharClass {
	return l.charClasses
}

func (l *MudLib) AddCharClass(newClass CharClass) {
	_, exists := l.charClasses[newClass.Name()]
	if !exists {
		l.charClasses[newClass.Name()] = newClass
	}
}

func (l *MudLib) RemoveCharClass(toRemove CharClass) bool {
	_, exists := l.charClasses[toRemove.Name()]
	delete(l.charClasses, toRemove.Name())
	return exists
}

func (l *MudLib) FindCharClass(search string) CharClass {
	// find exact match
	foundClass, exists := l.charClasses[search]
	if exists {
		return foundClass
	}

	search = strings.ToLower(search)
	// find partial match
	for name, charClass := range l.CharClasses() {
		if strings.HasPrefix(strings.ToLower(name), search) {
			return charClass
		}
	}

	return nil
}

func (l *MudLib) LoadCharClasses() {
	Fighter := PlayerClass{
		name:    "Fighter",
		realm:   RealmImmortal,
		enabled: false,
		statBonuses: [6]float64{
			1.0, 1.0, 1.0, 1.0, 1.0, 1.0,
		},
	}
	l.AddCharClass(&Fighter)

	Skeleton := PlayerClass{
		name:    "Skeleton",
		realm:   RealmEvil,
		enabled: true,
		statBonuses: [6]float64{
			1.5, 0.6, 1.4, 1.6, 0.4, 0.5,
		},
	}
	l.AddCharClass(&Skeleton)

	Necromancer := PlayerClass{
		name:    "Necromancer",
		realm:   RealmEvil,
		enabled: true,
		statBonuses: [6]float64{
			0.5, 0.6, 1.3, 0.7, 1.5, 1.4,
		},
	}
	l.AddCharClass(&Necromancer)
}

func (l *MudLib) Races() map[string]Race {
	return l.races
}

func (l *MudLib) AddRace(newRace Race) {
	_, exists := l.races[newRace.Name()]
	if !exists {
		l.races[newRace.Name()] = newRace
	}
}

func (l *MudLib) RemoveRace(toRemove Race) bool {
	_, exists := l.races[toRemove.Name()]
	delete(l.races, toRemove.Name())
	return exists
}

func (l *MudLib) FindRace(search string) Race {
	// find exact match
	race, exists := l.races[search]
	if exists {
		return race
	}

	search = strings.ToLower(search)
	// find partial match
	for name, pRace := range l.Races() {
		if strings.HasPrefix(strings.ToLower(name), search) {
			return pRace
		}
	}

	return nil
}

func (l *MudLib) LoadRaces() {
	Human := PlayerRace{
		name:    "Human",
		enabled: false,
		statBonuses: [6]float64{
			1.0, 1.0, 1.0, 1.0, 1.0, 1.0,
		},
	}
	l.AddRace(&Human)

	Skeleton := PlayerRace{
		name:    "Skeleton",
		enabled: true,
		statBonuses: [6]float64{
			1.5, 0.6, 1.4, 1.6, 0.4, 0.5,
		},
	}
	l.AddRace(&Skeleton)

	Necromancer := PlayerRace{
		name:    "Necromancer",
		enabled: true,
		statBonuses: [6]float64{
			0.5, 0.6, 1.3, 0.7, 1.5, 1.4,
		},
	}
	l.AddRace(&Necromancer)
}

func (l *MudLib) LoadCommands() {
	SayCmd := Command{
		trigger:    "say",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&SayCmd)

	LookCmd := Command{
		trigger:    "look",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&LookCmd)

	RenameCmd := Command{
		trigger:    "rename",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&RenameCmd)

	RerollCmd := Command{
		trigger:    "reroll",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&RerollCmd)

	SpawnCmd := Command{
		trigger:    "*spawn",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
		security:   SecAdmin,
	}
	l.AddCommand(&SpawnCmd)

	StatsCmd := Command{
		trigger:    "stats",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&StatsCmd)

	HealthCmd := Command{
		trigger:    "health",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&HealthCmd)

	ReleaseCmd := Command{
		trigger:    "release",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&ReleaseCmd)

	AttackCmd := Command{
		trigger:    "hit",
		timer:      2000,
		timerType:  TimerAttack,
		costType:   CostFat,
		useCost:    3,
		checkTimer: true,
	}
	l.AddCommand(&AttackCmd)

	InfoCmd := Command{
		trigger:    "information",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&InfoCmd)

	QuitCmd := Command{
		trigger:    "quit",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&QuitCmd)

	InventoryCmd := Command{
		trigger:    "inventory",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&InventoryCmd)

	WealthCmd := Command{
		trigger:    "wealth",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&WealthCmd)

	CreateCmd := Command{
		trigger:    "create",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
		security:   SecAdmin,
	}
	l.AddCommand(&CreateCmd)

	GetCmd := Command{
		trigger:    "get",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&GetCmd)

	DropCmd := Command{
		trigger:    "drop",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: true,
	}
	l.AddCommand(&DropCmd)

	GiveCmd := Command{
		trigger:    "give",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
	}
	l.AddCommand(&GiveCmd)

	GoCmd := Command{
		trigger:    "go",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: true,
	}
	l.AddCommand(&GoCmd)

	GoToCmd := Command{
		trigger:    "*goto",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: true,
		security:   SecStaff,
	}
	l.AddCommand(&GoToCmd)

	ShutdownCmd := Command{
		trigger:    "*shutdown",
		timer:      0,
		timerType:  TimerNone,
		costType:   CostNone,
		useCost:    0,
		checkTimer: false,
		security:   SecAdmin,
	}
	l.AddCommand(&ShutdownCmd)
}

func (l *MudLib) Combat() string {
	return l.combat
}
