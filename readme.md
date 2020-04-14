### Push Down Atomata in Go

#### Objectives

The objective of this assignment is to provide REST APIs to process the input tokens and to verify them against a PushDown Automata. The program has the following APIs:
- /pdas - returns id and name of all pdas
- /pdas/id - creates a new PDA from the request body provided in the input
- /pdas/id/reset - returns the PDA i.e. the stack and the queue of unprocessed tokens.
- /pdas/id/tokens/position - Process the token for the specified location. The token is provided as a part of the request body. If the token position is not one to be immediately processed, then the token is queued for later processing
- /pdas/id/eos/position - Specifies the position of the last token for the input. 
- /pdas/id/is_accepted - Specifies if the PDA is in accepting state or not.
- /pdas/id/stack/top/k - returns the top k symbols from the stack
- /pdas/id/stack/len - return the length of the stack
- /pdas/id/state - returns id and name of all pdas
- /pdas/id/tokens - returns the list of tokens that have been queued for later processing
- /pdas/id/snapshot/k - returns the current state, the queued tokens and the top k symbols in the stack.
- /pdas/id/close - 
- /pdas/id/delete - deletes the pda for the given id


#### How to run the program?

This is a Golang program that is run from the command line. 

Server startup Command: bash start_server.sh

**To create a PDA**

curl -X PUT -H "Content-Type: application/json" -d '{
    "name": "0n1n",
    "states": ["q1", "q2", "q3", "q4"],
    "input_alphabet": [ "0", "1" ],
    "stack_alphabet" : [ "0", "1" ],
    "accepting_states": ["q1", "q4"],
    "start_state": "q1",
    "transitions": [
        ["q1", "null", "null", "q2", "$"],
        ["q2", "0", "null", "q2", "0"],
        ["q2", "0", "0", "q2", "0"],
        ["q2", "1", "0", "q3", "null"],
        ["q3", "1", "0", "q3", "null"],
        ["q3", "null", "$", "q4", "null"]
    ],
    "eos": "$"
}' http://localhost:8080/pdas/100


**To get back all pdas**

curl -X GET http://localhost:8080/pdas


**To tokens in a PDA**

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/100/tokens/0

**To get current state of PDA**

curl -X GET http://localhost:8080/pdas/100/state


**To get Tokens**

curl -X GET http://localhost:8080/pdas/100/tokens

**To get snapshots**

curl -X GET http://localhost:8080/pdas/100/snapshot/3

**To get EOS**

curl http://localhost:8080/pdas/100/eos/6

**To Get is accepted**

curl http://localhost:8080/pdas/100/is_accepted