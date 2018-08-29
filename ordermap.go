package ordermap

import(
	"fmt"
	"../def"
	"time"
	"../hardware/elevator_io"
	"math"
)


/*-------------ORDERMAP------------------//
The ordermap is a three-dimensional matrix, 
represeting the orders of each elevator. 
This helps the elevators to keep track of everyones 
orders, both cab- and hall orders. 
Underneath is one third of an ordermap, showing which
order-information of each elevator is being stored. 

				up:	  down: cab:
Floor 1:		0     0   	0 

Floor 2:		0     0     0 

Floor 3:		0  	  0     0 

Floor 4:		0	  0     0

//---------------------------------------*/



var Current_floor                 int 
var Current_dir                   int 
var Is_elev_alive[def.NUM_ELEV]   bool
var Current_order_floor           int 
var Order_map[def.NUM_ELEV][def.NUM_FLOORS][def.NUM_BUTTONS] int 
var Elev_info_matrix[def.NUM_ELEV]Elev_info 


type Elev_info struct{
	Elev_ID                     int
	Current_dir                 int
	Current_floor               int
	Is_elev_alive[def.NUM_ELEV] bool
	Current_order_floor         int 
	Order_map[def.NUM_ELEV][def.NUM_FLOORS][def.NUM_BUTTONS] int 
}


func Ordermap_add_order_to_map(ID int, floor int, button int) bool {
	if Order_map[ID][floor][button] == 0 {	
		for elev := 0; elev < def.NUM_ELEV; elev++ {
			if button != 2 {
				Order_map[elev][floor][button] = 1 
			} 
		}
		if button == 2 { 
			Order_map[ID][floor][button] = 1
		}
		return true 
	}
	return false
}


func Ordermap_delete_order_from_map(f int, ID int){
	for elev := 0; elev < def.NUM_ELEV; elev++ {
		for button := 0; button < def.NUM_BUTTONS-1; button++ { 
			Order_map[elev][f][button] = 0			
		}  
	}
	if Order_map[ID-1][f][2] == 1 && Current_floor==f {
		elevator_io.SetButtonLamp(elevator_io.ButtonType(2), f, false)
	}
	if Order_map[ID-1][f][2] == 1 {
		Order_map[ID-1][f][2] = 0
	}
}


func Ordermap_choose_order_to_complete(set_dir chan int){
	for {
		if Elev_info_matrix[def.MyID-1].Current_order_floor == -1 { 
			for floor := 0; floor < def.NUM_FLOORS; floor++ {
				for button := 0; button < def.NUM_BUTTONS; button++{
					if Elev_info_matrix[def.MyID-1].Order_map[def.MyID-1][floor][button] == 1 { 
						set_dir <- floor
						break
					} 
				}
			}
		}
		time.Sleep(500 * time.Millisecond) 
	}
}


func Ordermap_find_dir_for_elev(floor int) elevator_io.MotorDirection {
	if floor > Elev_info_matrix[def.MyID-1].Current_floor {
		return elevator_io.MD_Up
	} 
	if floor < Elev_info_matrix[def.MyID-1].Current_floor {
		return elevator_io.MD_Down
	}
	return elevator_io.MD_Stop 
}


func Ordermap_am_I_closest(floor int ) bool { 
	result := true
	my_distance := int(math.Abs(float64(Elev_info_matrix[def.MyID-1].Current_floor - floor)))

	if Elev_info_matrix[def.MyID-1].Current_floor < floor { // Orders above
		for elev := 0; elev < def.NUM_ELEV; elev++{
			if (elev != def.MyID-1) && (Elev_info_matrix[def.MyID-1].Is_elev_alive[elev] == true) {
				elev_distance := int(math.Abs(float64(Elev_info_matrix[elev].Current_floor - floor)))

				if elev_distance < my_distance {
					return false

				} else if elev_distance == my_distance { 
					if Elev_info_matrix[elev].Elev_ID < Elev_info_matrix[def.MyID-1].Elev_ID {
						result = false  
					}
				}

			}
		}

	} else if Elev_info_matrix[def.MyID-1].Current_floor > floor { // Orders belove
		for elev := 0; elev < def.NUM_ELEV; elev++{
			if (elev != def.MyID-1) && (Elev_info_matrix[def.MyID-1].Is_elev_alive[elev] == true) { 
				elev_distance := int(math.Abs(float64(Elev_info_matrix[elev].Current_floor - floor)))

				if elev_distance < my_distance {
					return false 

				} else if elev_distance == my_distance {
					if Elev_info_matrix[elev].Elev_ID < Elev_info_matrix[def.MyID-1].Elev_ID {
						result = false  
					}
				}
			}
		}

	}
	return result 	
}


func Ordermap_check_if_order_at_current_floor(f int) bool {
	if Elev_info_matrix[def.MyID-1].Current_dir == 1 {
		if Elev_info_matrix[def.MyID-1].Order_map[def.MyID-1][f][0] == 1 || Elev_info_matrix[def.MyID-1].Order_map[def.MyID-1][f][2] == 1{
			return true 
		}
	}	
	if Elev_info_matrix[def.MyID-1].Current_dir == -1 {
		if Elev_info_matrix[def.MyID-1].Order_map[def.MyID-1][f][1] == 1 || Elev_info_matrix[def.MyID-1].Order_map[def.MyID-1][f][2] == 1 {
			return true 
		}
	}	
	return false 
}


func Ordermap_print_order_map(){
	fmt.Println(" ** ORDER MAP **", "\n")
	for elev := 0; elev < def.NUM_ELEV; elev++ {
		for floor := 0; floor < def.NUM_FLOORS; floor++ {
			fmt.Println("Elev", elev, ":", Order_map[elev][floor])
		}
		fmt.Println("\n")
	}
}