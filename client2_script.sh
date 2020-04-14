printf "\n---------------- Create PDA with id 101 -----------------" 
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
}' http://localhost:8080/pdas/101

printf "\n PDA created"

curl -X GET http://localhost:8080/pdas


printf "\n------------Put tokens ---------------\n" 

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/101/tokens/0

printf "\n"

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/101/tokens/1


printf "\n---------------- Current state of the PDA ------------------\n" 

curl -X GET http://localhost:8080/pdas/101/state

printf "\n---------------- Current length of stack -------------------\n" 

curl -X GET http://localhost:8080/pdas/101/stack/len

printf "\n-------------- Continue processing other tokens ----------\n" 

curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/101/tokens/2

printf "\n"

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/101/tokens/3


printf "\n---------------- Snapshot ----------------------------------\n" 

curl -X GET http://localhost:8080/pdas/101/snapshot/3

printf "\n---------------- Call eos ----------------------------------\n" 

curl http://localhost:8080/pdas/101/eos/6

printf "\n---------------- Call is_accepted() -------------------\n" 

curl http://localhost:8080/pdas/101/is_accepted


printf "\n---------------- Reset ----------------------------------\n" 

curl -X GET http://localhost:8080/pdas/101/reset

printf "\n---------------- Snapshot ----------------------------------\n" 

curl -X GET http://localhost:8080/pdas/101/snapshot/3

printf "\n---------------- Show all Pdas -------------------\n" 

curl -X GET http://localhost:8080/pdas

printf "\n---------------- Delete Pdas -------------------\n" 

curl -X GET http://localhost:8080/pdas/101/delete

printf "\n---------------- Show all Pdas -------------------\n" 

curl -X GET http://localhost:8080/pdas

