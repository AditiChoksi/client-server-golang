package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
)

// A structure that defines the a Push Down Automata
// This class contains the attributes of a PDA
type Pda struct {
	Id string
	Name string
	States [] string
	Input_alphabet [] string
	Stack_alphabet [] string
	Accepting_states [] string
	Start_state string
	Transitions [][]string
	Eos string
}

// This class simulates a PDA processor that is it runs the PDA for teh provided input.
type PDAProcessor struct{
	Stack [] string
	Pda Pda
	Current_State string
}

var pdaList []Pda

func createPda(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Pdas Endpoint Hit")
	var p Pda

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}
	
	// Do something with the Pda struct...
	fmt.Fprintf(w, "Pda: %+v", p)
	pdaList = []Pda {p}
	
	json.NewEncoder(w).Encode(pdaList)

}

func returnAllPdas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Pdas Endpoint Hit")
	json.NewEncoder(w).Encode(pdaList)
}


func  handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/pdas", returnAllPdas)
	myRouter.HandleFunc("/pdas/{id}", createPda)
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main(){
	fmt.Println("Hi")

	handleRequest()
}