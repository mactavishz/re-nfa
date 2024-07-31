package fsm

import "unicode/utf8"

type State struct {
	transitions map[rune][]*State
	epsilon     []*State
	isAccepting bool
}

type NFA struct {
	start  *State
	accept *State
}

func (s *State) AddTransition(to *State, input rune) {
	s.transitions[input] = append(s.transitions[input], to)
}

func (s *State) AddEpsilon(to *State) {
	s.epsilon = append(s.epsilon, to)
}

func NewNFA(start *State, accept *State) *NFA {
	return &NFA{start: start, accept: accept}
}

func NewState(isAccepting bool) *State {
	return &State{
		transitions: make(map[rune][]*State),
		epsilon:     []*State{},
		isAccepting: isAccepting,
	}
}

// catenation operator
func Concat(left, right *NFA) *NFA {
	left.accept.isAccepting = false
	left.accept.AddEpsilon(right.start)
	return &NFA{start: left.start, accept: right.accept}
}

// '|' operator
func Or(left, right *NFA) *NFA {
	start := NewState(false)
	start.AddEpsilon(left.start)
	start.AddEpsilon(right.start)

	accept := NewState(true)
	left.accept.isAccepting = false
	right.accept.isAccepting = false
	left.accept.AddEpsilon(accept)
	right.accept.AddEpsilon(accept)

	return &NFA{start: start, accept: accept}
}

// '*' operator
func Star(nfa *NFA) *NFA {
	start := NewState(false)
	accept := NewState(true)

	start.AddEpsilon(nfa.start)
	start.AddEpsilon(accept)

	nfa.accept.isAccepting = false
	nfa.accept.AddEpsilon(nfa.start)
	nfa.accept.AddEpsilon(accept)

	return &NFA{start: start, accept: accept}
}

// '?' operator
func Optional(nfa *NFA) *NFA {
	start := NewState(false)
	accept := NewState(true)

	start.AddEpsilon(nfa.start)
	start.AddEpsilon(accept)

	nfa.accept.AddEpsilon(accept)
	nfa.accept.isAccepting = false

	return &NFA{start: start, accept: accept}
}

// '+' operator
func Plus(nfa *NFA) *NFA {
	start := NewState(false)
	accept := NewState(true)

	start.AddEpsilon(nfa.start)
	nfa.accept.AddEpsilon(nfa.start)
	nfa.accept.AddEpsilon(accept)
	nfa.accept.isAccepting = false

	return &NFA{start: start, accept: accept}
}

func (nfa *NFA) Match(input string) bool {
	current := map[*State]bool{nfa.start: true}

	// process a potential chain of epsilon transitions
	processEpsilon := func() {
		changed := true
		for changed {
			changed = false
			for state := range current {
				for _, eps := range state.epsilon {
					if !current[eps] {
						current[eps] = true
						changed = true
					}
				}
			}
		}
	}

	processEpsilon()

	// Check if the empty string is accepted
	if len(input) == 0 {
		for state := range current {
			if state.isAccepting {
				return true
			}
		}
		return false
	}

	// Process each character in the input
	for len(input) > 0 {
		ch, size := utf8.DecodeRuneInString(input)
		input = input[size:]

		next := make(map[*State]bool)
		for state := range current {
			if nextStates, ok := state.transitions[ch]; ok {
				for _, nextState := range nextStates {
					next[nextState] = true
				}
			}
		}

		current = next
		processEpsilon()
	}

	// Only check for accepting state after processing all input
	for state := range current {
		if state.isAccepting {
			return true
		}
	}

	return false
}
