echo "-----------------------Client1-----------------" 
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

curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/100/tokens/1
curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/100/tokens/2
curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/100/tokens/3
curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/100/tokens/4
curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/100/tokens/5
curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/100/tokens/0
curl http://localhost:8080/pdas/100/eos/6
curl http://localhost:8080/pdas/100/is_accepted

echo "-----------------------Client2-----------------" 
curl -X PUT -H "Content-Type: application/json" -d '{
    "name": "1n0n",
    "states": ["q1", "q2", "q3", "q4"],
    "input_alphabet": [ "0", "1" ],
    "stack_alphabet" : [ "0", "1" ],
    "accepting_states": ["q1", "q4"],
    "start_state": "q1",
    "transitions": [
        ["q1", "null", "null", "q2", "$"],
        ["q2", "1", "null", "q2", "1"],
        ["q2", "1", "1", "q2", "1"],
        ["q2", "0", "1", "q3", "null"],
        ["q3", "0", "1", "q3", "null"],
        ["q3", "null", "$", "q4", "null"]
    ],
    "eos": "$"
}' http://localhost:8080/pdas/101


curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/101/tokens/1
curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/101/tokens/2
curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/101/tokens/3
curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/101/tokens/4
curl -X PUT -H "Content-Type: application/json" -d '{"token": "0"}' http://localhost:8080/pdas/101/tokens/5
curl -X PUT -H "Content-Type: application/json" -d '{"token": "1"}' http://localhost:8080/pdas/101/tokens/0



curl http://localhost:8080/pdas/101/eos/6
curl http://localhost:8080/pdas/101/is_accepted
