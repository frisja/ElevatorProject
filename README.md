# TTK4145 - Elevator project
 
Our design is based on peer to peer principal, where all elevators know everything about each other. The elevators uses this information to decide if they are closest to the current order. The peers module, which was handed out, is being used to keep track of elevators on the network ( and elevators that are supposed to be on the network). 

The main tasks of the different modules are as follows:

-ordermap: It’s main purpose is to keep track of the orders for every elevator. It is also keeping track of all the information of every elevator.  

-def: this module is housing global constants and variables being used by other modules. 

-timer: the timer is keeping track of time being used serving an order.

-network: sends information of the elevators over the internet. 

-hardware: is the interface between hardware and software.

-fsm: in this module the states of an elevator is being changed. 

We are not authors of bcast, peers, conn and elevator_io (although they have been edited).

