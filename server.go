package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"time"
	"strconv"
)

// This class simulates a PDA processor that is it runs the PDA for teh provided input.
type PDAProcessor struct{
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
}

//var pdaList []PDAProcessor
var c = cache.New(5*time.Minute, 10*time.Minute)

// Function to push data on to the stack when executing the PDA. It modifies the stack.
func push(p PDAProcessor, val string) {
	p.Stack = append(p.Stack, val)
}

// Function to pop data from the stack when executing the PDA. It modifies the stack.
func pop(p PDAProcessor) {
	p.Stack = p.Stack[:len(p.Stack) -1]
}

// Function to obtain the top n elements of the stack. This function does not modify the stack.
func peekInternal(p PDAProcessor, k int) []string {
	
	top := [] string{}
	l := len(p.Stack)
	if (l <= k) {
		top = p.Stack
	} else if ( k == 1) {
		top = append(top, p.Stack[l-1])
	} else {
		top = p.Stack[l-k:l-1]
	}
	return top
}

// Function to pop data from the stack when executing the PDA. It modifies the stack.
func stacklen(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]

	var p, _ = c.Get(id)
	proc := p.(*PDAProcessor)

	var l = len(proc.Stack)

	json.NewEncoder(w).Encode(l)
}

// Function to obtain the top n elements of the stack. This function does not modify the stack.
func peek(w http.ResponseWriter, r *http.Request) {

	var vars = mux.Vars(r)
	var id = vars["id"]
	var kstring = vars["k"]
	k, _ := strconv.Atoi(kstring)

	var p, _ = c.Get(id)
	proc := p.(PDAProcessor)

	top := peekInternal(proc, k)
	
	json.NewEncoder(w).Encode(top)
	//return top
}

// API to reset the PDA and the stack. This deletes everything from the stack 
// and sets the current state to the start state so that we can start new.
func reset(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var pda, _ = c.Get(id)
	p := pda.(PDAProcessor)
	p.Stack = make([]string, 0)
	p.Current_State = p.Start_state
}

func createPda(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Create Pdas")
	var p PDAProcessor
	
	var vars = mux.Vars(r)
	var id = vars["id"]

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}
	
	c.Set(id, p, cache.NoExpiration)
	//fmt.Fprintf(w, "cache: %+v", c)
	//var x, _ = c.Get(id)
	//fmt.Fprintf(w, "Pda: %+v", x)
	
	json.NewEncoder(w).Encode(c.Items())

}

func returnAllPdas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Return all Pdas")
	json.NewEncoder(w).Encode(c.Items())
}


// Function to open the grammar file. This function unmarshal's the data input into the PDA structure.
func open(data []byte, p *PDAProcessor) bool {	
	err := json.Unmarshal(data, &p)
	if err != nil {
		fmt.Print(err)
		return false
	}
	return true
}

// Function to check if the input string has been accepted by the pda 
func is_accepted(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]

	var p, _ = c.Get(id)
	proc := p.(*PDAProcessor)

	flag := false
	accepting_states := proc.Accepting_states
	cs := proc.Current_State

	if len(proc.Stack) == 0 {
		for i:= 0; i < len(accepting_states); i++ {
			if cs == accepting_states[i] {
				flag = true
				fmt.Println("\n***************************")
				fmt.Println("Input token Accepted.")
				break
			}
		}
	}
	if !flag {
		fmt.Println("\n***************************")
		fmt.Println("Input string Rejected.")
	}

	json.NewEncoder(w).Encode(flag)
	//return flag
}

// The done returns the final status of the current state and the stack after the input string is processed.
func done(proc PDAProcessor, is_accepted bool, transition_count int) {
	fmt.Println("pda = ", proc.Name,"::total_clock = ", transition_count, "::method = is_accepted = ", is_accepted,"::Current State = ", proc.Current_State)
	fmt.Println("Current_state: ", proc.Current_State)
	fmt.Println(proc.Stack)
}

// Returns the current state of the PDA
func current_state(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]

	var p, _ = c.Get(id)
	proc := p.(*PDAProcessor)

	json.NewEncoder(w).Encode(proc.Current_State)

}

func put(w http.ResponseWriter, r *http.Request) {
fmt.Println("Endpoint Hit: Put")
	var vars = mux.Vars(r)
	var id = vars["id"]
	var token = vars["token"]

	var p, _ = c.Get(id)
	proc := p.(PDAProcessor)
	var transition_count = putInternal(proc,token)

	json.NewEncoder(w).Encode(transition_count)
}

// This function accepts the input string and performs the necessary transitions and 
// stack operations for every token,
func putInternal(proc PDAProcessor, token string) int{
	

	transitions := proc.Transitions
	tran_len := len(transitions)
	transition_count := 0
	for j := 0; j < tran_len; j++ {
		var allowed_current_state = transitions[j][0]
		var input = transitions[j][1]
		var allowed_top_of_stack = transitions[j][2]
		var target_state = transitions[j][3]
		var action_item = transitions[j][4]
		var currentStackSymbol = ""
		var top = peekInternal(proc, 1)
		if(len(top)>=1){
			currentStackSymbol = top[0]
		}
		
		// PDA is deterministic. It jumps from current state to target state in the specified conditions
		if (input == "null" && allowed_current_state == proc.Current_State && allowed_top_of_stack == "null" && action_item == "null") {
			fmt.Println("Current State ",proc.Current_State)
			fmt.Println("No push/pop performed...... Processed dead transition")
			fmt.Println("Stack: ", proc.Stack)
			fmt.Println("New State ", target_state)
			proc.Current_State = target_state
			transition_count = transition_count + 1
		}

		if (allowed_current_state == proc.Current_State && input == token)  {

			//Perform Push action
			if action_item != "null" && allowed_top_of_stack == "null" {
				fmt.Println("Current State ", proc.Current_State)
				fmt.Println("Push ", action_item, " on the stack.")
				fmt.Println("New State ", target_state)
				fmt.Println("Stack: ", proc.Stack)
				transition_count = transition_count + 1
				proc.Current_State = target_state
				push(proc, action_item)

				break
				//performs Push action
			} else if action_item != "null" &&  allowed_top_of_stack == currentStackSymbol {
				fmt.Println("Current State ",proc.Current_State)
				fmt.Println("Push ", action_item, " on the stack")
				fmt.Println("New State ", target_state)
				fmt.Println("Stack: ", proc.Stack)
				transition_count = transition_count + 1
				proc.Current_State = target_state
				push(proc, action_item)
				break
				//performs Pop action
			} else if action_item == "null" &&  allowed_top_of_stack == currentStackSymbol {
				pop(proc)
				fmt.Println("Current State ",proc.Current_State)
				fmt.Println("Pop top of the stack.")
				fmt.Println("New State ",target_state)
				fmt.Println("Stack: ", proc.Stack)
				transition_count = transition_count + 1
				proc.Current_State = target_state
				break
				//Neither push nor pop action required
			} else if allowed_top_of_stack == "null" {
				fmt.Println("Current State ",proc.Current_State)
				fmt.Println("No push/pop performed...... Consumed input token")
				fmt.Println("New State ",target_state)
				fmt.Println("Stack: ", proc.Stack)
				proc.Current_State = target_state
				transition_count = transition_count + 1
				break
			}
		}	       
	}

	//c.Set(id, proc, cache.NoExpiration)

	fmt.Println("Clock count for consuming the input token = ", transition_count)
	//json.NewEncoder(w).Encode(transition_count)
	return transition_count
}

// Performs the last transition to move the Automata to accepting state after the input
// string has been successfully parsed. 
func eos(w http.ResponseWriter, r *http.Request) {

	var vars = mux.Vars(r)
	var id = vars["id"]

	var p, _ = c.Get(id)
	proc := p.(PDAProcessor)

	length_of_stack := len(proc.Stack)
	allowed_transitions := proc.Transitions
	target_state := ""
	allowed_top_of_stack := ""
	var currentStackSymbol = ""
	var top = peekInternal(proc, 1)
	if(len(top)>=1){
		currentStackSymbol = top[0]
	}
	for j := 0; j < len(allowed_transitions); j++ {	
		var allowed_current_state = allowed_transitions[j][0]
		allowed_top_of_stack = allowed_transitions[j][2]
		
		if allowed_current_state == proc.Current_State && allowed_top_of_stack == currentStackSymbol{
			target_state = allowed_transitions[j][3]
			break
		}
	}
	if currentStackSymbol == proc.Eos {
		fmt.Println("")
		fmt.Println("Popping last $ from the stack")
		fmt.Println("Current State ",proc.Current_State)
		fmt.Println("New State ",target_state)
		proc.Current_State = target_state
		if length_of_stack > 0 {
			pop(proc)
		}
	}
}

// Pushes initial EOS token into the stack and moves to the next state indicating
// the start of transitions
func check_for_first_move(proc PDAProcessor, transition_count int) int{
	allowed_transitions := proc.Transitions
	target_state := ""
	input := ""
	allowed_top_of_stack := ""
	action_item := ""
	
	for j := 0; j < len(allowed_transitions); j++ {
		var allowed_current_state = allowed_transitions[j][0]
		if allowed_current_state == proc.Current_State {
			input = allowed_transitions[j][1]
			allowed_top_of_stack = allowed_transitions[j][2]
			target_state = allowed_transitions[j][3]
			action_item = allowed_transitions[j][4]
			break
		}
	}
	
	if input == "null" && allowed_top_of_stack == "null"{
		fmt.Println("Current State ", proc.Current_State)

		push(proc, action_item)
		fmt.Println("Pushing $ on the stack")

		proc.Current_State = target_state
		fmt.Println("New State ", proc.Current_State)
        
		transition_count = transition_count + 1
		fmt.Println()
	} 
	return transition_count
}

//Checks whether the input string is composed of the allowed characters. 
func verify_Input_String(proc PDAProcessor, input_string string)bool{
	var input_symbols = proc.Input_alphabet
	verify:=false
	for i :=0; i < len(input_string); i++ {
		verify=false
		for j :=0; j < len(input_symbols); j++ {
			if string(input_string[i]) == input_symbols[j] {
				verify = true
				break
			}
		}
		
		if verify == false {
			break
		}
	}
	return verify
}

func  handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/pdas", returnAllPdas)
	myRouter.HandleFunc("/pdas/{id}", createPda)
	myRouter.HandleFunc("/pdas/{id}/reset", reset)
	myRouter.HandleFunc("/pdas/{id}/{token}/{position}", put)
	myRouter.HandleFunc("/pdas/{id}/eos/{position}", eos)
	myRouter.HandleFunc("/pdas/{id}/is_accepted", is_accepted)
	myRouter.HandleFunc("/pdas/{id}/stack/top/{k}", peek)
	myRouter.HandleFunc("/pdas/{id}/stack/len", stacklen)
	myRouter.HandleFunc("/pdas/{id}/state", current_state)
	//myRouter.HandleFunc("/pdas/{id}/tokens", gettokens)
	//myRouter.HandleFunc("/pdas/{id}/snapshot/{k}", snapshot)
	//myRouter.HandleFunc("/pdas/{id}/close", close)
	//myRouter.HandleFunc("/pdas/{id}/delete", delete)


	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main(){
	fmt.Println("Server started. Listening at port 8080")

	handleRequest()
}