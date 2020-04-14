package main

// This class simulates a PDA processor that is it runs the PDA for teh provided input.
type PDAProcessor struct{
	Id string
	Name string
	States [] string
	Input_alphabet [] string
	Stack_alphabet [] string
	Accepting_states [] string
	Start_state string
	Transitions [][]string
	Eos string
	Stack [] string
	Current_State string
	Next_Position int
	Hold_back_Queue [] HoldBackStruct
}

type PDAInfo struct {
	Id string
	Name string
}

type Snapshot struct { 
	Topk [] string
	Current_State string
	Hold_back_Queue [] HoldBackStruct
}

type Token struct {
	Token string
}

type HoldBackStruct struct {
	Hold_back_Position string
	Hold_back_Token string
}