# Library Management System
## BE WARY OF PORTS BEING USED
### API-GATEWAY
#### Port Number used: 8081

### Book Service
#### Port Number used: 8080

### Fine Service
#### Port Number used: 

### Borrowing Service
#### Port Number used: 8082

### Patron Service
#### Port Number used: 8069
<br>
Notes for Front-End:
Inig himo sa patron-creation page, applyi og regex ang phone number na field <br>
Mao ni ako gi gamit : '^[0-9]{10,15}$'
<br><br>
Notes for communicating with Patron-Service:<br>
Patron Service makes us of rabbitmq queues which are seperated into two main queues:<br>

1. "patron-service-queue"
    - this queue is used exclusively for API-GATEWAY
2. "patron-service-internal-queue" 
    - this queue is used exclusively for interservice communication
    - use this queue in order to access the mutations, queries, and subscriptions of the patron service



