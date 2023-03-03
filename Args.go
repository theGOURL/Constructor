package constructor

type Args interface {
	// Get retorna o enésimo argumento, ou uma string em branco
	Get(n int) string
	// First retorna o primeiro argumento, ou então uma string em branco
	First() string
	// Tail retorna o restante dos argumentos (não o primeiro)
	// ou senão uma fatia de string vazia
	Tail() []string
	// Len retorna o comprimento da fatia embrulhada
	Len() int
	// Apresenta verificações se existem argumentos presentes
	Present() bool
	// Slice retorna uma cópia da fatia interna
	Slice() []string
}

type args []string

func (a *args) Get(n int) string {
	if len(*a) > n {
		return (*a)[n]
	}
	return ""
}

func (a *args) First() string {
	return a.Get(0)
}

func (a *args) Tail() []string {
	if a.Len() >= 2 {
		tail := []string((*a)[1:])
		ret := make([]string, len(tail))
		copy(ret, tail)
		return ret
	}
	return []string{}
}

func (a *args) Len() int {
	return len(*a)
}

func (a *args) Present() bool {
	return a.Len() != 0
}

func (a *args) Slice() []string {
	ret := make([]string, len(*a))
	copy(ret, *a)
	return ret
}
