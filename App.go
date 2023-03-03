package constructor

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//ignoreFlagPrefix é  a constante que ignora sinalizadores de teste ao adicionar sinalizadores de outros pacotes
//ignoreFlagPrefix é  a constante que ignora sinalizadores de teste ao adicionar sinalizadores de outros pacotes
const ignoreFlagPrefix = "test."

//App é a struct Principal para cria uma aplicação constructor
//App é a struct Principal para cria uma aplicação constructor
type App struct {
	//O nome do programa. O padrão é path.Base(os.Args[0])
//O nome do programa. O padrão é path.Base(os.Args[0])
	Name string
	//O nome completo do comando de help.
//O nome completo do comando de help.
	HelpName string
	//descriçãoDaconstructor
//descriçãoDaconstructor
	Usage string
	//Texto para substituir a seção USAGE de help
//Texto para substituir a seção USAGE de help
	UsageText string
	//Descrição do formato do argumento do programa.
//Descrição do formato do argumento do programa.
	ArgsUsage string
	//Versão do programa
//Versão do programa
	Version string
	//Descrição do programa
//Descrição do programa
	Description string
	//DefaultCommand é o nome (opcional) de um comando
//DefaultCommand é o nome (opcional) de um comando
	//para executar se nenhum nome de comando for passado como argumentos constructor.
//para executar se nenhum nome de comando for passado como argumentos constructor.
	DefaultCommand string
	//Lista de Comandos para executar
//Lista de Comandos para executar
	Commands []*Command
	//Lista de sinalizadores para analisar
//Lista de sinalizadores para analisar
	Flags []Flag
	//Conclusão da concha de Booleano para Habilatarataratar
//Conclusão da concha de Booleano para Habilatarataratar
	EnableShellCompletion bool
	//nomeDoComandoDeGeraçãoDeConclusãoDeShell
//nomeDoComandoDeGeraçãoDeConclusãoDeShell
	ShellCompletionCommandName string
	//Booleano para ocultar o comando de ajuda integrado e o sinalizador de ajuda
//Booleano para ocultar o comando de ajuda integrado e o sinalizador de ajuda
	HideHelp bool
	//Booleano para ocultar o comando de ajuda integrado, mas manter o sinalizador de ajuda.
//Booleano para ocultar o comando de ajuda integrado, mas manter o sinalizador de ajuda.
	//Ignorado se HideHelp for verdadeiro.
//Ignorado se HideHelp for verdadeiro.
	HideHelpCommand bool
	//Booleano para ocultar o sinalizador de versão integrado e a seção VERSION da ajuda
//Booleano para ocultar o sinalizador de versão integrado e a seção VERSION da ajuda
	HideVersion bool
	//categorias contém os comandos categorizados e é preenchido na inicialização do aplicativo
//categorias contém os comandos categorizados e é preenchido na inicialização do aplicativo
	categories CommandCategories
	// flagCategories contains the categorized flags and is populated on app startup
// FlagCategories contém as bandeiras categorizadas e é preenchido na inicialização do aplicativo
	flagCategories FlagCategories
	// An action to execute when the shell completion flag is set
// Uma ação a ser executada quando o sinalizador de conclusão do shell é definido
	ShellComplete ShellCompleteFunc
	// An action to execute before any subcommands are run, but after the context is ready
// uma ação a ser executada antes que qualquer subcompôs seja executado, mas depois que o contexto estiver pronto
	// If a non-nil error is returned, no subcommands are run
// Se um erro não-nil for retornado, nenhum subcomando é executado
	Before BeforeFunc
	// An action to execute after any subcommands are run, but after the subcommand has finished
// uma ação a ser executada após a execução de quaisquer subcompâncias forem executados, mas depois que o subcomando ter terminado
	// It is run even if Action() panics
// é executado mesmo se Action () pânico
	After AfterFunc
	// The action to execute when no subcommands are specified
// a ação a ser executada quando nenhum subcompôs é especificado
	Action ActionFunc
	// Execute this function if the proper command cannot be found
// executar esta função se o comando adequado não puder ser encontrado
	CommandNotFound CommandNotFoundFunc
	// Execute this function if a usage error occurs
// execute esta função se ocorrer um erro de uso
	OnUsageError OnUsageErrorFunc
	// Execute this function when an invalid flag is accessed from the context
// executa esta função quando uma bandeira inválida é acessada do contexto
	InvalidFlagAccessHandler InvalidFlagAccessFunc
	// List of all authors who contributed (string or fmt.Stringer)
// Lista de todos os autores que contribuíram (string ou fmt.stringer)
	Authors []any // TODO: ~string | fmt.Stringer when interface unions are available
 // TODO: ~ String |fmt.stringer quando os sindicatos da interface estão disponíveis
	// Copyright of the binary if any
// direitos autorais do binário, se houver
	Copyright string
	// Reader reader to write input to (useful for tests)
// leitor leitor para escrever a entrada (útil para testes)
	Reader io.Reader
	// Writer writer to write output to
// escritor escritor para escrever saída para
	Writer io.Writer
	// ErrWriter writes error output
// errwriter grava saída de erro
	ErrWriter io.Writer
	// ExitErrHandler processes any error encountered while running an App before
// O ExiterrHandler processa qualquer erro encontrado ao executar um aplicativo antes
	// it is returned to the caller. If no function is provided, HandleExitCoder
// é devolvido ao chamador.Se nenhuma função for fornecida, HandleExitCoder
	// is used as the default behavior.
// é usado como comportamento padrão.
	ExitErrHandler ExitErrHandlerFunc
	// Other custom info
// Outras informações personalizadas
// retornará mais cedo se a configuração já tivesse acontecido.
	Metadata map[string]interface{}
	// Carries a function which returns app specific info.
// carrega uma função que retorna informações específicas do aplicativo.
	ExtraInfo func() map[string]string
	// CustomAppHelpTemplate the text template for app help topic.
// CustomAppHelPtemplate O modelo de texto para tópico de ajuda do aplicativo.
	// constructor.go uses text/template to render templates. You can
// constructor.go usa texto/modelo para renderizar modelos.Você pode
	// render custom help text by setting this variable.
// renderiza o texto de ajuda personalizado configurando essa variável.
	CustomAppHelpTemplate string
	// SliceFlagSeparator is used to customize the separator for SliceFlag, the default is ","
// sliceflagseparator é usado para personalizar o separador para Sliceflag, o padrão é ","
	SliceFlagSeparator string
	// DisableSliceFlagSeparator is used to disable SliceFlagSeparator, the default is false
// desabilitaliceflagseparator é usado para desativar o sliceflagseparator, o padrão é falso
	DisableSliceFlagSeparator bool
	// Boolean to enable short-option handling so user can combine several
// booleano para permitir o manuseio de opção curta para que o usuário possa combinar vários
	// single-character bool arguments into one
// Argumentos de bool
	// i.e. foobar -o -v -> foobar -ov
// ou seja,Fubar -o -in -> Fubar -oo
	UseShortOptionHandling bool
	// Enable suggestions for commands and flags
// Ative sugestões para comandos e bandeiras
	Suggest bool
	// Allows global flags set by libraries which use flag.XXXVar(...) directly
// permite sinalizadores globais definidos por bibliotecas que usam flag.xxxvar (...) diretamente
	// to be parsed through this library
// para ser analisado através desta biblioteca
	AllowExtFlags bool
	// Treat all flags as normal arguments if true
// trate todas as bandeiras como argumentos normais se verdadeiro
	SkipFlagParsing bool
	// Flag exclusion group
// Grupo de exclusão de bandeira
	MutuallyExclusiveFlags []MutuallyExclusiveFlags
	// Use longest prefix match for commands
// Use a correspondência de prefixo mais longa para comandos
	PrefixMatchCommands bool
	// Custom suggest command for matching
// Sugira comando de sugestão personalizado para correspondência
	SuggestCommandFunc SuggestCommandFunc

	didSetup bool

	rootCommand *Command

	// if the app is in error mode
// se o aplicativo estiver no modo de erro
	isInError bool
}

// Setup runs initialization code to ensure all data structures are ready for
// A configuração executa o código de inicialização para garantir que todas as estruturas de dados estejam prontas para
// `Run` or inspection prior to `Run`.  It is internally called by `Run`, but
// `run` ou inspeção antes do` run`.É chamado internamente por `run`, mas
// will return early if setup has already happened.
// retornará mais cedo se a configuração já tivesse acontecido.
func (a *App) Setup() {
	if a.didSetup {
		return
	}

	a.didSetup = true

	if a.Name == "" {
		a.Name = filepath.Base(os.Args[0])
	}

	if a.HelpName == "" {
		a.HelpName = a.Name
	}

	if a.Usage == "" {
		a.Usage = "A new constructor application"
	}

	if a.Version == "" {
		a.HideVersion = true
	}

	if a.ShellComplete == nil {
		a.ShellComplete = DefaultAppComplete
	}

	if a.Action == nil {
		a.Action = helpCommand.Action
	}

	if a.Reader == nil {
		a.Reader = os.Stdin
	}

	if a.Writer == nil {
		a.Writer = os.Stdout
	}

	if a.ErrWriter == nil {
		a.ErrWriter = os.Stderr
	}

	if a.AllowExtFlags {
		// Adicione bandeiras globais adicionadas por outros pacotespor outros pacotesor outros pacotesor outros pacotes
// Adicione bandeiras globais adicionadas por outros pacotespor outros pacotesor outros pacotesor outros pacotes
		flag.VisitAll(func(f *flag.Flag) {
			//Pule sinalizadores de testees de testees de teste
//Pule sinalizadores de testes de testes de teste
			if !strings.HasPrefix(f.Name, ignoreFlagPrefix) {
				a.Flags = append(a.Flags, &extFlag{f})
			}
		})
	}

	var newCommands []*Command

	for _, c := range a.Commands {
		cname := c.Name
		if c.HelpName != "" {
			cname = c.HelpName
		}
		c.HelpName = fmt.Sprintf("%s %s", a.HelpName, cname)

		c.flagCategories = newFlagCategoriesFromFlags(c.Flags)
		newCommands = append(newCommands, c)
	}
	a.Commands = newCommands

	if a.Command(helpCommand.Name) == nil && !a.HideHelp {
		if !a.HideHelpCommand {
			helpCommand.HelpName = fmt.Sprintf("%s %s", a.HelpName, helpName)
			a.appendCommand(helpCommand)
		}

		if HelpFlag != nil {
			a.appendFlag(HelpFlag)
		}
	}

	if !a.HideVersion {
		a.appendFlag(VersionFlag)
	}

	if a.PrefixMatchCommands {
		if a.SuggestCommandFunc == nil {
			a.SuggestCommandFunc = suggestCommand
		}
	}
	if a.EnableShellCompletion {
		if a.ShellCompletionCommandName != "" {
			completionCommand.Name = a.ShellCompletionCommandName
		}
		a.appendCommand(completionCommand)
	}

	a.categories = newCommandCategories()
	for _, command := range a.Commands {
		a.categories.AddCommand(command.Category, command)
	}
	sort.Sort(a.categories.(*commandCategories))

	a.flagCategories = newFlagCategories()
	for _, fl := range a.Flags {
		if cf, ok := fl.(CategorizableFlag); ok {
			if cf.GetCategory() != "" {
				a.flagCategories.AddFlag(cf.GetCategory(), cf)
			}
		}
	}

	if a.Metadata == nil {
		a.Metadata = make(map[string]interface{})
	}

	if len(a.SliceFlagSeparator) != 0 {
		defaultSliceFlagSeparator = a.SliceFlagSeparator
	}

	disableSliceFlagSeparator = a.DisableSliceFlagSeparator
}

func (a *App) newRootCommand() *Command {
	return &Command{
		Name:                   a.Name,
		Usage:                  a.Usage,
		UsageText:              a.UsageText,
		Description:            a.Description,
		ArgsUsage:              a.ArgsUsage,
		ShellComplete:          a.ShellComplete,
		Before:                 a.Before,
		After:                  a.After,
		Action:                 a.Action,
		OnUsageError:           a.OnUsageError,
		Commands:               a.Commands,
		Flags:                  a.Flags,
		flagCategories:         a.flagCategories,
		HideHelp:               a.HideHelp,
		HideHelpCommand:        a.HideHelpCommand,
		UseShortOptionHandling: a.UseShortOptionHandling,
		HelpName:               a.HelpName,
		CustomHelpTemplate:     a.CustomAppHelpTemplate,
		categories:             a.categories,
		SkipFlagParsing:        a.SkipFlagParsing,
		isRoot:                 true,
		MutuallyExclusiveFlags: a.MutuallyExclusiveFlags,
		PrefixMatchCommands:    a.PrefixMatchCommands,
	}
}

func (a *App) newFlagSet() (*flag.FlagSet, error) {
	return flagSet(a.Name, a.Flags)
}

func (a *App) useShortOptionHandling() bool {
	return a.UseShortOptionHandling
}

//Run é o ponto de entrada para o aplicativo constructor.Analisa os argumentos fatia e rotasrotasrotas
//Run é o ponto de entrada para o aplicativo constructor.Analisa os argumentos fatia e rotasrotasrotas
//para a combinação de bandeira/args adequadadequadadequada
//para a combinação de bandeira/args adequadadequadadequada
func (a *App) Run(arguments []string) (err error) {
	return a.RunContext(context.Background(), arguments)
}

//RunContext é como Run, exceto que é preciso um contexto que serárárá
//RunContext é como Run, exceto que é preciso um contexto que serárárá
//passou para seus comandos e subcomandos.Com isso, você pode
//passou para seus comandos e subcomandos.Com isso, você pode
//Propagar tempo limite e solicitações de cancelamentolamentolamento
//Propagar tempo limite e solicitações de cancelamentolamentolamento
func (a *App) RunContext(ctx context.Context, arguments []string) (err error) {
	a.Setup()

	//manuseie a bandeira de conclusão separadamente do flagset desdedede
//manuseie a bandeira de conclusão separadamente do flagset desdedede
	//A conclusão pode ser tentada após uma bandeira, mas antes que seu valor fosse colocadoosse colocadoosse colocado
//A conclusão pode ser tentada após uma bandeira, mas antes que seu valor fosse colocadoosse colocadoosse colocado
	//na linha de comando.Isso faz com que o flagset interprete a conclusão
//na linha de comando.Isso faz com que o flagset interprete a conclusão
	//Nome da bandeira como o valor da bandeira antes dele, o que é indesejávelsejávelsejável
//Nome da bandeira como o valor da bandeira antes dele, o que é indesejávelsejávelsejável
	//Observe que só podemos fazer isso porque a função de preenchimento automático da Shellomático da Shellomático da Shell
//Observe que só podemos fazer isso porque a função de preenchimento automático da Shellomático da Shellomático da Shell
	//sempre anexa a bandeira de conclusão no final do comando
//sempre anexa a bandeira de conclusão no final do comando
	shellComplete, arguments := checkShellCompleteFlag(a, arguments)

	cCtx := NewContext(a, nil, &Context{Context: ctx})
	cCtx.shellComplete = shellComplete

	a.rootCommand = a.newRootCommand()
	cCtx.Command = a.rootCommand

	return a.rootCommand.Run(cCtx, arguments...)
}

func (a *App) suggestFlagFromError(err error, command string) (string, error) {
	flag, parseErr := flagFromError(err)
	if parseErr != nil {
		return "", err
	}

	flags := a.Flags
	hideHelp := a.HideHelp
	if command != "" {
		cmd := a.Command(command)
		if cmd == nil {
			return "", err
		}
		flags = cmd.Flags
		hideHelp = hideHelp || cmd.HideHelp
	}

	suggestion := SuggestFlag(flags, flag, hideHelp)
	if len(suggestion) == 0 {
		return "", err
	}

	return fmt.Sprintf(SuggestDidYouMeanTemplate+"\n\n", suggestion), nil
}

//O comando retorna o comando nomeado no aplicativo.Retorna nil se o comando não existiririr
//O comando retorna o comando nomeado no aplicativo.Retorna nil se o comando não existiririr
func (a *App) Command(name string) *Command {
	for _, c := range a.Commands {
		if c.HasName(name) {
			return c
		}
	}

	return nil
}

//VisibleCategorias retorna uma fatia de categorias e comandos que são
//VisibleCategorias retorna uma fatia de categorias e comandos que são
//Oculto = falseee
// Hidden = falsee
func (a *App) VisibleCategories() []CommandCategory {
	ret := []CommandCategory{}
	for _, category := range a.categories.Categories() {
		if visible := func() CommandCategory {
			if len(category.VisibleCommands()) > 0 {
				return category
			}
			return nil
		}(); visible != nil {
			ret = append(ret, visible)
		}
	}
	return ret
}

//VisibleCommands retorna uma fatia dos comandos com Hidden = false
//VisibleCommands retorna uma fatia dos comandos com Hidden = false
func (a *App) VisibleCommands() []*Command {
	var ret []*Command
	for _, command := range a.Commands {
		if !command.Hidden {
			ret = append(ret, command)
		}
	}
	return ret
}

//VisibleFlagCategorias retorna uma fatia contendo todas as categorias com as bandeiras que eles contêmontêmontêm
//Disable Flag Categorias retorna uma fatia contendo todas as categorias com as bandeiras que eles contêmontêmontêm
func (a *App) VisibleFlagCategories() []VisibleFlagCategory {
	if a.flagCategories == nil {
		return []VisibleFlagCategory{}
	}
	return a.flagCategories.VisibleCategories()
}

//VisibleFlags retorna uma fatia das bandeiras com Hidden = falselselse
//VisibleFlags retorna uma fatia das bandeiras com Hidden = falselselse
func (a *App) VisibleFlags() []Flag {
	return visibleFlags(a.Flags)
}

func (a *App) appendFlag(fl Flag) {
	if !hasFlag(a.Flags, fl) {
		a.Flags = append(a.Flags, fl)
	}
}

func (a *App) appendCommand(c *Command) {
	if !hasCommand(a.Commands, c) {
		a.Commands = append(a.Commands, c)
	}
}

func (a *App) handleExitCoder(cCtx *Context, err error) {
	if a.ExitErrHandler != nil {
		a.ExitErrHandler(cCtx, err)
	} else {
		HandleExitCoder(err)
	}
}

func (a *App) argsWithDefaultCommand(oldArgs Args) Args {
	if a.DefaultCommand != "" {
		rawArgs := append([]string{a.DefaultCommand}, oldArgs.Slice()...)
		newArgs := args(rawArgs)

		return &newArgs
	}

	return oldArgs
}

func runFlagActions(c *Context, fs []Flag) error {
	for _, f := range fs {
		isSet := false
		for _, name := range f.Names() {
			if c.IsSet(name) {
				isSet = true
				break
			}
		}
		if isSet {
			if af, ok := f.(ActionableFlag); ok {
				if err := af.RunAction(c); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a *App) writer() io.Writer {
	if a.isInError {
		//Isso pode acontecer em teste, mas não em uso normalmalmal
//Isso pode acontecer em teste, mas não em uso normalmalmal
		if a.ErrWriter == nil {
			return os.Stderr
		}
		return a.ErrWriter
	}
	return a.Writer
}

func checkStringSliceIncludes(want string, sSlice []string) bool {
	found := false
	for _, s := range sSlice {
		if want == s {
			found = true
			break
		}
	}

	return found
}
