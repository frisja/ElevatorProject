package fsm 

import(
	"../ordermap" 
	"../hardware/elevator_io"
	"../def"
	"time"
	"../network/network"
	"sync"
	"../timer"

)

var mutex = &sync.Mutex{}

func Fsm_set_dir_for_elev(set_dir chan int, delete_order_trans chan network.Delete_order_from_map, get_current_floor chan int, peer_tx_enable chan bool) {
	for{
		select{
			case floor :=<- set_dir: 
				if floor != -1 {
					is_closest := ordermap.Ordermap_am_I_closest(floor)

					if (is_closest || ordermap.Order_map[def.MyID-1][floor][2] == 1) { // hall order || cab order 
						if ordermap.Current_order_floor == -1 { // If the elevator has no current order 
							ordermap.Current_order_floor = floor
							dir := ordermap.Ordermap_find_dir_for_elev(floor)	
							
							mutex.Lock()
							ordermap.Current_dir = int(dir)
							elevator_io.SetMotorDirection(dir)
							mutex.Unlock()
							
							go timer.Timer_is_elev_completing_order(get_current_floor,set_dir, peer_tx_enable) 
							
							if dir == 0 { // already here 
								Fsm_door_open()
								network.Network_order_complete(floor, delete_order_trans)
							}
						}
					}
				}
		}
		time.Sleep(500 * time.Millisecond)
	}
}


func Fsm_door_open(){
	elevator_io.SetDoorOpenLamp(true)
  	time.Sleep(def.TIME_DOOR_OPEN*time.Second)
	elevator_io.SetDoorOpenLamp(false)
}


func Fsm_arrived_on_ordered_floor(reached_floor int, delete_order_trans chan network.Delete_order_from_map) bool{  
	if reached_floor == ordermap.Elev_info_matrix[def.MyID-1].Current_order_floor {
		return true
	}

	if ordermap.Ordermap_check_if_order_at_current_floor(reached_floor) { 
		elevator_io.SetMotorDirection(elevator_io.MD_Stop)
		Fsm_door_open() 
		network.Network_order_complete(reached_floor, delete_order_trans)
	} 
	
	return false	
}

func Fsm_emergency_stop() { 
	elevator_io.SetMotorDirection(elevator_io.MD_Stop)
	elevator_io.SetStopLamp(true)

	if elevator_io.GetFloor() != -1 { 
    	Fsm_door_open()
    }

    time.Sleep(def.EMERGENCY_STOP_TIME * time.Second)
    elevator_io.SetStopLamp(false)

    if ordermap.Current_order_floor!=-1 {
    	dir := ordermap.Ordermap_find_dir_for_elev(ordermap.Current_order_floor)
    	mutex.Lock()
		ordermap.Current_dir = int(dir)
		elevator_io.SetMotorDirection(dir)
		mutex.Unlock()   
    }
   	
	 
}

func Fsm_obstruction() {
	elevator_io.SetMotorDirection(elevator_io.MD_Stop)
	for elevator_io.GetObstruction(){
		// While the obstruction is active 
	}

	for {  
		if elevator_io.GetFloor()==int(-1){
			elevator_io.SetMotorDirection(elevator_io.MD_Up)
		} 
		if elevator_io.GetFloor()!=int(-1){
			elevator_io.SetMotorDirection(elevator_io.MD_Stop)
			break
		}
		
	}
}
