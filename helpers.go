package main

import (
	"fmt"
	"strconv"
)

// Function to push data on to the stack when executing the PDA. It modifies the stack.
func push(p *PDAProcessor, val string) {
	p.Stack = append(p.Stack, val)
}

// Function to pop data from the stack when executing the PDA. It modifies the stack.
func pop(p *PDAProcessor) {
	p.Stack = p.Stack[:len(p.Stack) -1]
}

// Function to obtain the top n elements of the stack. This function does not modify the stack.
func peekInternal(p *PDAProcessor, k int) []string {
	
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

// API to reset the PDA and the stack. This deletes everything from the stack 
// and sets the current state to the start state so that we can start new.
func resetInternal(p *PDAProcessor) {
	p.Stack = make([]string, 0)
	p.Current_State = p.Start_state
	p.Next_Position = 0
	p.Hold_back_Queue = make([]HoldBackStruct, 0)
}

// create the PDA struct from the request json
func open(id string, p PDAProcessor) bool {
	
	p.Id = id

	_, found := cache[id]

	if !found {
		resetInternal(&p)
		cache[id] = p
		return true
	}

	return false
}

// Function to check if the input string has been accepted by the pda 
func is_accepted_internal(proc PDAProcessor) bool{
	flag := false
	accepting_states := proc.Accepting_states
	cs := proc.Current_State
	if len(proc.Stack) == 0 && len(proc.Hold_back_Queue) == 0{
		for i:= 0; i < len(accepting_states); i++ {
			if cs == accepting_states[i] {
				flag = true
				break
			}
		}
	}
	return flag
}

func process_hold_back_tokens(proc PDAProcessor) int{
	defer wg.Done()
	token_blocked := -1
	token_processed := false
	for {
		if(len(proc.Hold_back_Queue) == 0) {
			break
		}
		proc = cache[proc.Id]
		hold_back := proc.Hold_back_Queue[len(proc.Hold_back_Queue) -1]
		pos_int, _ := strconv.Atoi(hold_back.Hold_back_Position)
		if(proc.Next_Position == pos_int) {
			token_processed = putInternal(proc,hold_back.Hold_back_Token)
			if(!token_processed) {
				token_blocked = pos_int
				break
			} else
			{
				proc = cache[proc.Id]
				proc.Hold_back_Queue = proc.Hold_back_Queue[:len(proc.Hold_back_Queue) -1]
				cache[proc.Id] = proc

			}
		} else 
		{
			break
		}
    	
	}
	return token_blocked	
}

// This function accepts the input string and performs the necessary transitions and 
// stack operations for every token,
func putInternal(proc PDAProcessor, token string) bool{	
	transitions := proc.Transitions
	tran_len := len(transitions)
	token_processed := false
	for j := 0; j < tran_len; j++ {
		var allowed_current_state = transitions[j][0]
		var input = transitions[j][1]
		var allowed_top_of_stack = transitions[j][2]
		var target_state = transitions[j][3]
		var action_item = transitions[j][4]
		var currentStackSymbol = ""
		var top = peekInternal(&proc, 1)
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
		}

		if (allowed_current_state == proc.Current_State && input == token)  {

			//Perform Push action
			if action_item != "null" && allowed_top_of_stack == "null" {
				fmt.Println("Current State ", proc.Current_State)
				fmt.Println("Push ", action_item, " on the stack.")
				fmt.Println("New State ", target_state)
				fmt.Println("Stack: ", proc.Stack)
				proc.Next_Position = proc.Next_Position + 1
				token_processed = true
				proc.Current_State = target_state
				push(&proc, action_item)

				break
				//performs Push action
			} else if action_item != "null" &&  allowed_top_of_stack == currentStackSymbol {
				fmt.Println("Current State ",proc.Current_State)
				fmt.Println("Push ", action_item, " on the stack")
				fmt.Println("New State ", target_state)
				fmt.Println("Stack: ", proc.Stack)
				proc.Next_Position = proc.Next_Position + 1
				token_processed = true
				proc.Current_State = target_state
				push(&proc, action_item)
				break
				//performs Pop action
			} else if action_item == "null" &&  allowed_top_of_stack == currentStackSymbol {
				pop(&proc)
				fmt.Println("Current State ",proc.Current_State)
				fmt.Println("Pop top of the stack.")
				fmt.Println("New State ",target_state)
				fmt.Println("Stack: ", proc.Stack)
				proc.Next_Position = proc.Next_Position + 1
				token_processed = true
				proc.Current_State = target_state
				break
				//Neither push nor pop action required
			} else if allowed_top_of_stack == "null" {
				fmt.Println("Current State ",proc.Current_State)
				fmt.Println("No push/pop performed...... Consumed input token")
				fmt.Println("New State ",target_state)
				fmt.Println("Stack: ", proc.Stack)
				proc.Current_State = target_state
				proc.Next_Position = proc.Next_Position + 1
				token_processed = true
				break
			}
		}	       
	}
	cache[proc.Id] = proc	
	return token_processed
}

// Pushes initial EOS token into the stack and moves to the next state indicating
// the start of transitions
func check_for_first_move(proc *PDAProcessor, transition_count int){
	transitions := proc.Transitions
	target_state := ""
	input := ""
	allowed_top_of_stack := ""
	action_item := ""
	
	for j := 0; j < len(transitions); j++ {

		if transitions[j][0] == proc.Current_State {
			input = transitions[j][1]
			allowed_top_of_stack = transitions[j][2]
			target_state = transitions[j][3]
			action_item = transitions[j][4]
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

}