package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
)

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
}

var pdaList []PDAProcessor

// Function to push data on to the stack when executing the PDA. It modifies the stack.
func push(p *PDAProcessor, val string) {
	p.Stack = append(p.Stack, val)
}

// Function to pop data from the stack when executing the PDA. It modifies the stack.
func pop(p *PDAProcessor) {
	p.Stack = p.Stack[:len(p.Stack) -1]
}

// Function to obtain the top n elements of the stack. This function does not modify the stack.
func peek(p *PDAProcessor, k int) []string {
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
func reset(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	for i := 0; i < len(pdaList); i++ {
		if pdaList[i].Id == id {
			pdaList[i].Stack = make([]string, 0)
			pdaList[i].Current_State = pdaList[i].Start_state
		}
	}
}

func createPda(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Create Pdas Endpoint Hit")
	var p PDAProcessor

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}
	
	// Do something with the Pda struct...
	fmt.Fprintf(w, "Pda: %+v", p)
	pdaList = []PDAProcessor {p}
	
	json.NewEncoder(w).Encode(pdaList)

}

func returnAllPdas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Return all Pdas Endpoint Hit")
	json.NewEncoder(w).Encode(pdaList)
}


func  handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/pdas", returnAllPdas)
	myRouter.HandleFunc("/pdas/{id}", createPda)
	//myRouter.HandleFunc("/pdas/{id}", createPda)
	myRouter.HandleFunc("/pdas/{id}/reset", reset)
	// myRouter.HandleFunc("/pdas/{id}/is_accepted", reset)
	// myRouter.HandleFunc("/pdas/{id}/reset", reset)


	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main(){
	fmt.Println("Hi")

	handleRequest()
}