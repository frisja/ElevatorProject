package network 

import(
	"time"
	"../../ordermap"
	"../../hardware/elevator_io"
	"../peers"
	"../../def"
	"strconv"
	"sync"	
)

var mutex = &sync.Mutex{}

type Button_pressed_message struct{	
	Elev_ID            int
	Floor 			   int
	Button 			   int
}	 


type Take_order_msg struct {
	Floor              int 
	Direction          int 
	Elev_take_order_ID int
}


type Delete_order_from_map struct{
	Floor              int
	Elev_ID            int
}


func Network_update_info(info_trans_chan chan ordermap.Elev_info){  
	for{
		info_msg := ordermap.Elev_info{def.MyID, ordermap.Current_dir, ordermap.Current_floor, ordermap.Is_elev_alive, ordermap.Current_order_floor, ordermap.Order_map}

		mutex.Lock()
		ordermap.Elev_info_matrix[def.MyID-1] = info_msg 
		mutex.Unlock()

    	info_trans_chan <- info_msg	
    	time.Sleep(50 * time.Millisecond) 
    }
}


func Network_order_complete(floor int, trans_delete_order chan Delete_order_from_map){
	ordermap.Ordermap_delete_order_from_map(floor, def.MyID)

	mutex.Lock()
	ordermap.Current_order_floor = -1
	ordermap.Current_dir = 0
	mutex.Unlock()

	delete_msg := Delete_order_from_map{floor, def.MyID}

	for rm_msg_iteration := 0; rm_msg_iteration < 100; rm_msg_iteration++ { // to ensure that the message is received
		trans_delete_order <- delete_msg
	}
}				



func Network_button_pushed(floor int, button int, add_order_trans chan Button_pressed_message){	
	
	change_made := ordermap.Ordermap_add_order_to_map(def.MyID-1, floor, button)

	if change_made { 
		msg := Button_pressed_message{def.MyID-1, floor, int(button)} 
		
		for add_msg := 0; add_msg < 100; add_msg++ { 
			add_order_trans <- msg
		}

	} 

	if (button == 2 || button == 5 || button == 8) && ordermap.Elev_info_matrix[def.Elev_alive-1].Is_elev_alive[def.Elev_alive-1]==true {
		elevator_io.SetButtonLamp(elevator_io.ButtonType(button), floor, true)
	}
	
}


func Network_elev_alive_control(peer_update chan peers.PeerUpdate){
	for{
		select{
			case update :=<-peer_update:

				for elev := 0; elev < def.NUM_ELEV; elev++ {
					mutex.Lock()
					ordermap.Is_elev_alive[elev] = false  
					mutex.Unlock()
				}
	
				for elev := 0; elev < len(update.Peers); elev++{
					elev_int, _ := strconv.Atoi(update.Peers[elev])
					for my_id := 1; my_id < 4; my_id++ {
						if elev_int == my_id {
							mutex.Lock()
							ordermap.Is_elev_alive[my_id-1] = true 
							mutex.Unlock()
							break
						}
					}
				}

				if len(update.Peers) > 0 {
					first_peer_in_list, _ := strconv.Atoi(update.Peers[0])
					def.Elev_alive = first_peer_in_list
					def.Num_elev_alive = len(update.Peers)
				}

				if len(update.Peers) ==2 {
					first_peer_in_list, _ := strconv.Atoi(update.Peers[0])
					second_peer_in_list, _ := strconv.Atoi(update.Peers[1])
					def.Elev_alive1 =first_peer_in_list
					def.Elev_alive2 =second_peer_in_list
				}


				var new_elev int
				var port string 

				if len(update.New) > 0 {  
					for elev := 0; elev < len(update.New); elev++ {
						new_elev = int(update.New[elev])- 48

						if new_elev == 1 {
							port = ":15657"
						}						
						if new_elev == 2 {
							port = ":15658"	
						}				
						if new_elev == 3 {
							port = ":15659"	
						}				
					}

					elevator_io.Init(port, def.NUM_FLOORS)
					if  elevator_io.GetFloor()!=-1{
						ordermap.Current_floor=elevator_io.GetFloor()
					}
					

					if len(update.Peers) > 0 { 
						if len(update.Peers) > 1 && def.MyID == 1 {
							elev_int, _ := strconv.Atoi(update.Peers[1])
							ordermap.Order_map = ordermap.Elev_info_matrix[elev_int-1].Order_map 					
						}else if len(update.Peers) > 1 {
							elev_peers1, _ := strconv.Atoi(update.Peers[0])
							if elev_peers1==2 && def.MyID==2 {
								elev_int, _ := strconv.Atoi(update.Peers[1])
								ordermap.Order_map = ordermap.Elev_info_matrix[elev_int-1].Order_map 
							}					
						} else{
							ordermap.Order_map = ordermap.Elev_info_matrix[def.Elev_alive-1].Order_map 
						}		
					}
				}	
		}
		time.Sleep(100 * time.Millisecond)
	}
}



