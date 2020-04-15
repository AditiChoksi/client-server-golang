package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"strconv"
	"sort"
	"sync"
)

var wg sync.WaitGroup
var cache = make(map[string]PDAProcessor) 

// Function to pop data from the stack when executing the PDA. It modifies the stack.
func stacklen(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]

	proc := cache[id]
	var l = len(proc.Stack)

	json.NewEncoder(w).Encode(l)
}

// Function to obtain the top n elements of the stack. This function does not modify the stack.
func peek(w http.ResponseWriter, r *http.Request) {

	var vars = mux.Vars(r)
	var id = vars["id"]
	var kstring = vars["k"]
	k, _ := strconv.Atoi(kstring)

	proc := cache[id]
	top := peekInternal(&proc, k)
	
	json.NewEncoder(w).Encode(top)

}

// API to reset the PDA and the stack. This deletes everything from the stack 
// and sets the current state to the start state so that we can start new.
func reset(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]

	p := cache[id]
	resetInternal(&p)
	cache[id] = p
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
	created := open(id, p)

	if created {
		json.NewEncoder(w).Encode("PDA successfully created.")
	} else 
	{
		json.NewEncoder(w).Encode("Cannot create PDA. A PDA with this id already exists.")
	}
}

func returnAllPdas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Return all Pdas")
	var pdalist[]PDAInfo
	for key, value := range cache {
		info := PDAInfo{
			Id: key,
			Name: value.Name,
		}
		pdalist = append(pdalist, info)
	}
	json.NewEncoder(w).Encode(pdalist)
}

// Function to check if the input string has been accepted by the pda 
func is_accepted(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]

	proc := cache[id]

	flag := is_accepted_internal(proc)
	if(flag){
		json.NewEncoder(w).Encode("Input tokens successfully Accepted")
	} else
	{
		json.NewEncoder(w).Encode("Input tokens Rejected by the PDA")
	}
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

	proc := cache[id]

	json.NewEncoder(w).Encode(proc.Current_State)

}

func put(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Put")
	var vars = mux.Vars(r)
	var id = vars["id"]
	var position = vars["position"]

	var t Token
	json.NewDecoder(r.Body).Decode(&t)
	var token = t.Token

	proc := cache[id]
	check_for_first_move(&proc, 1)
	pos_int, _ := strconv.Atoi(position)

	token_processed := false
	token_blocked := -1
	var hold_back_flag = false
	if(proc.Next_Position == pos_int) {
		fmt.Println ("Calling Put")

		token_processed = putInternal(proc,token)
		if(token_processed) {
			wg.Add(1)
			go func() {
   				token_blocked = process_hold_back_tokens(proc)
			}()
			wg.Wait()
		} else
		{
			token_blocked = pos_int
		}

	} else if (proc.Next_Position < pos_int) {
		var duplicate_token = false


		for _, v := range proc.Hold_back_Queue {
			hold_back_pos_int, _ := strconv.Atoi(v.Hold_back_Position)
    		if hold_back_pos_int == pos_int {
        		duplicate_token = true
    		}
		}
		if(!duplicate_token) {
			var hold_back HoldBackStruct
			hold_back_flag = true
			hold_back.Hold_back_Token = token
			hold_back.Hold_back_Position = position

			proc.Hold_back_Queue = append(proc.Hold_back_Queue , hold_back)
			sort.Slice(proc.Hold_back_Queue, func(i, j int) bool {
				return proc.Hold_back_Queue[i].Hold_back_Position > proc.Hold_back_Queue[j].Hold_back_Position
			})
			cache[proc.Id] = proc
			json.NewEncoder(w).Encode("Token kept in hold_back_Queue")
		} else 
		{
			json.NewEncoder(w).Encode("Duplicate token received")
		}

	} else 
	{
		hold_back_flag = true
		json.NewEncoder(w).Encode("Conflicting values for token at the same position.")
	}	
	if (token_blocked == -1 && !hold_back_flag) {
		json.NewEncoder(w).Encode("Token processed successfully. Please enter the next token")
	} else if (!hold_back_flag) {
		flag := is_accepted_internal(proc)
		if(flag){
			json.NewEncoder(w).Encode("Input tokens successfully Accepted")
		} else
		{
			json.NewEncoder(w).Encode("Input tokens Rejected by the PDA")
		}
	} 
}



func gettokens(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Get Queued tokens")
	var vars = mux.Vars(r)
	var id = vars["id"]

	proc := cache[id]

	for j := 0; j < len(proc.Hold_back_Queue)-1; j++ {
		fmt.Println("Queued token :", proc.Hold_back_Queue[j].Hold_back_Token, " At position :", proc.Hold_back_Queue[j].Hold_back_Position)
	}
	json.NewEncoder(w).Encode(proc.Hold_back_Queue)
}

func snapshot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Get snapshot")
	var vars = mux.Vars(r)
	var id = vars["id"]

	proc := cache[id]
	var snap Snapshot
	snap.Topk = make([]string, 0)
	snap.Current_State = proc.Current_State
	snap.Hold_back_Queue = proc.Hold_back_Queue
	snap.Topk = peekInternal(&proc,5)
	json.NewEncoder(w).Encode(snap)
}

// Performs the last transition to move the Automata to accepting state after the input
// string has been successfully parsed. 
func eos(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: Eos")

	var vars = mux.Vars(r)
	var id = vars["id"]
	var position = vars["position"]
	pos_int, _ := strconv.Atoi(position)

	proc := cache[id]

	length_of_stack := len(proc.Stack)
	transitions := proc.Transitions
	target_state := ""
	allowed_top_of_stack := ""
	var currentStackSymbol = ""
	var top = peekInternal(&proc, 1)
	
	if(len(top)>=1){
		currentStackSymbol = top[0]
	}
	for j := 0; j < len(transitions); j++ {	
		var allowed_current_state = transitions[j][0]
		allowed_top_of_stack = transitions[j][2]
		
		if allowed_current_state == proc.Current_State && allowed_top_of_stack == currentStackSymbol{
			target_state = transitions[j][3]
			break
		}
	}
	if currentStackSymbol == proc.Eos && pos_int == proc.Next_Position{
		fmt.Println("")
		fmt.Println("Popping last $ from the stack")
		fmt.Println("Current State ",proc.Current_State)
		fmt.Println("New State ",target_state)
		proc.Current_State = target_state
		if length_of_stack > 0 {
			pop(&proc)
		}
	}

	cache[id] = proc
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

func deletePda(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]

	_, found := cache[id]

	if found {
		delete(cache, id)
		json.NewEncoder(w).Encode("Pda deleted.")
	} else {
		json.NewEncoder(w).Encode("Pda not found.")

	}
}

func close(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Success. No resources to clean.")
}

func  handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/pdas", returnAllPdas)
	myRouter.HandleFunc("/pdas/{id}", createPda)
	myRouter.HandleFunc("/pdas/{id}/reset", reset)
	myRouter.HandleFunc("/pdas/{id}/tokens/{position}", put)
	myRouter.HandleFunc("/pdas/{id}/eos/{position}", eos)
	myRouter.HandleFunc("/pdas/{id}/is_accepted", is_accepted)
	myRouter.HandleFunc("/pdas/{id}/stack/top/{k}", peek)
	myRouter.HandleFunc("/pdas/{id}/stack/len", stacklen)
	myRouter.HandleFunc("/pdas/{id}/state", current_state)
	myRouter.HandleFunc("/pdas/{id}/tokens", gettokens)
	myRouter.HandleFunc("/pdas/{id}/snapshot/{k}", snapshot)
	myRouter.HandleFunc("/pdas/{id}/close", close)
	myRouter.HandleFunc("/pdas/{id}/delete", deletePda)


	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main(){
	fmt.Println("Server started. Listening at port 8080")

	handleRequest()
}