package constructor

import "sort"

// CommandCategories é uma interface de categorias de comando permite manipulação de categoria
type CommandCategories interface {
	// AddCommand adiciona um comando a uma categoria, criando uma nova categoria, se necessário.
	AddCommand(category string, command *Command)
	// Categories retorna uma fatia de categorias classificadas pelo nome
	Categories() []CommandCategory
}

type commandCategories []*commandCategory

func newCommandCategories() CommandCategories {
	ret := commandCategories([]*commandCategory{})
	return &ret
}

func (c *commandCategories) Less(i, j int) bool {
	return lexicographicLess((*c)[i].Name(), (*c)[j].Name())
}

func (c *commandCategories) Len() int {
	return len(*c)
}

func (c *commandCategories) Swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}

func (c *commandCategories) AddCommand(category string, command *Command) {
	for _, commandCategory := range []*commandCategory(*c) {
		if commandCategory.name == category {
			commandCategory.commands = append(commandCategory.commands, command)
			return
		}
	}
	newVal := append(*c,
		&commandCategory{name: category, commands: []*Command{command}})
	*c = newVal
}

func (c *commandCategories) Categories() []CommandCategory {
	ret := make([]CommandCategory, len(*c))
	for i, cat := range *c {
		ret[i] = cat
	}
	return ret
}

// CommandCategory é uma categoria que contém comandos.
type CommandCategory interface {
	// Name retorna a string de nome da categoria
	Name() string
	// VisibleCommands retorna uma fatia dos comandos com hidden = false
	VisibleCommands() []*Command
}

type commandCategory struct {
	name     string
	commands []*Command
}

func (c *commandCategory) Name() string {
	return c.name
}

func (c *commandCategory) VisibleCommands() []*Command {
	if c.commands == nil {
		c.commands = []*Command{}
	}

	var ret []*Command
	for _, command := range c.commands {
		if !command.Hidden {
			ret = append(ret, command)
		}
	}
	return ret
}

// A interface das categorias de sinalização permite a manipulação da categoria
type FlagCategories interface {
	// AddFlags adiciona um sinalizador a uma categoria, criando uma nova categoria, se necessário.
	AddFlag(category string, fl Flag)
	// VisibleCategorias retorna uma fatia de categorias de bandeira visível classificadas pelo nome
	VisibleCategories() []VisibleFlagCategory
}

type defaultFlagCategories struct {
	m map[string]*defaultVisibleFlagCategory
}

func newFlagCategories() FlagCategories {
	return &defaultFlagCategories{
		m: map[string]*defaultVisibleFlagCategory{},
	}
}

func newFlagCategoriesFromFlags(fs []Flag) FlagCategories {
	fc := newFlagCategories()
	for _, fl := range fs {
		if cf, ok := fl.(CategorizableFlag); ok {
			if cf.GetCategory() != "" {
				fc.AddFlag(cf.GetCategory(), cf)
			}
		}
	}

	return fc
}

func (f *defaultFlagCategories) AddFlag(category string, fl Flag) {
	if _, ok := f.m[category]; !ok {
		f.m[category] = &defaultVisibleFlagCategory{name: category, m: map[string]Flag{}}
	}

	f.m[category].m[fl.String()] = fl
}

func (f *defaultFlagCategories) VisibleCategories() []VisibleFlagCategory {
	catNames := []string{}
	for name := range f.m {
		catNames = append(catNames, name)
	}

	sort.Strings(catNames)

	ret := make([]VisibleFlagCategory, len(catNames))
	for i, name := range catNames {
		ret[i] = f.m[name]
	}

	return ret
}

// VisibleFlagCategory é uma categoria que contém sinalizadores.
type VisibleFlagCategory interface {
	// Name retorna a string de nome da categoria
	Name() string
	// Flags retorna uma fatia de visível flag classificada pelo nome
	Flags() []VisibleFlag
}

type defaultVisibleFlagCategory struct {
	name string
	m    map[string]Flag
}

func (fc *defaultVisibleFlagCategory) Name() string {
	return fc.name
}

func (fc *defaultVisibleFlagCategory) Flags() []VisibleFlag {
	vfNames := []string{}
	for flName, fl := range fc.m {
		if vf, ok := fl.(VisibleFlag); ok {
			if vf.IsVisible() {
				vfNames = append(vfNames, flName)
			}
		}
	}

	sort.Strings(vfNames)

	ret := make([]VisibleFlag, len(vfNames))
	for i, flName := range vfNames {
		ret[i] = fc.m[flName].(VisibleFlag)
	}

	return ret
}
