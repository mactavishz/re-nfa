package nfa

type Graph struct {
	label string
	start *State
	end   *State
}

type State struct {
	label string
	ch    rune
	out1  *State
	out2  *State
}

type SplitState struct {
	State
	isSplit bool
}

type AcceptState struct {
	State
	isAccept bool
}
