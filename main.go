package main

/*This is the entry point for the elevator project in TTK4145 Real time programming.
The project consists of six modules tied together in this main package. The modules communicate through go channels.
The modules are: def, fsm, hardware, network, ordermap and timer. Their communication and
further function description can be found in the README file.*/

import (
	"./ordermap"
	"./hardware/elevator_io"
	"./fsm"
	"./network/network"
	"./network/bcast"
	"sync"
	"./hardware/hardware"
	"flag"
	"./network/peers"	
	"./timer"
)

var mutex = &sync.Mutex{}

func main(){

	var id string 
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	peer_update          := make(chan peers.PeerUpdate)
	peer_tx_enable       := make(chan bool)
	 
	add_order_trans      := make(chan network.Button_pressed_message, 100) 
	add_order_rec        := make(chan network.Button_pressed_message, 100) 
	
	elev_info_trans      := make(chan ordermap.Elev_info, 100)
	elev_info_rec        := make(chan ordermap.Elev_info, 100)
	
	delete_order_trans   := make(chan network.Delete_order_from_map, 100)
	delete_order_rec     := make(chan network.Delete_order_from_map, 100)
	
	obstruction          := make(chan bool)
	emergency_stop       := make(chan bool)
	new_floor            := make(chan int)
	new_order            := make(chan elevator_io.ButtonEvent)
	set_dir 	         := make(chan int)
	get_current_floor    := make(chan int)


	hardware.Hardware_init(new_floor, new_order)

	go peers.Peers_transmitter(20020, id, peer_tx_enable)
	go peers.Peers_receiver(20020, peer_update)

	go bcast.Bcast_transmitter(20005, add_order_trans, elev_info_trans, delete_order_trans)
	go bcast.Bcast_receiver(20005, add_order_rec, elev_info_rec, delete_order_rec)
	
	go elevator_io.Elevator_io_poll_obstruction_switch(obstruction)
	go elevator_io.Elevator_io_poll_stop_button(emergency_stop)

	go network.Network_update_info(elev_info_trans)	
	go network.Network_elev_alive_control(peer_update)
	
	go timer.Timer_get_current_floor(get_current_floor)

	go fsm.Fsm_set_dir_for_elev(set_dir, delete_order_trans, get_current_floor, peer_tx_enable)

	go ordermap.Ordermap_choose_order_to_complete(set_dir)

	

	for{
		select{
		case newfloor :=<- new_floor:
			peer_tx_enable <- true
			ordermap.Current_floor = newfloor
			elevator_io.SetFloorIndicator(newfloor)
			
			should_stop := fsm.Fsm_arrived_on_ordered_floor(newfloor, delete_order_trans)

			if should_stop {
				elevator_io.SetMotorDirection(elevator_io.MD_Stop)
				fsm.Fsm_door_open() 
				network.Network_order_complete(newfloor, delete_order_trans)
			}


		case neworder :=<- new_order:
			network.Network_button_pushed(neworder.Floor, int(neworder.Button), add_order_trans)
			

		case recinfo :=<- elev_info_rec:
			mutex.Lock()
			ordermap.Elev_info_matrix[recinfo.Elev_ID-1] = recinfo 
			mutex.Unlock()


		case recmap :=<- add_order_rec: 
			ordermap.Ordermap_add_order_to_map(recmap.Elev_ID, recmap.Floor, recmap.Button)			


		case deletefloor :=<- delete_order_rec:
			ordermap.Ordermap_delete_order_from_map(deletefloor.Floor, deletefloor.Elev_ID)


		case <- obstruction:
			fsm.Fsm_obstruction()


		case <- emergency_stop:
			fsm.Fsm_emergency_stop()
		}
	}
}