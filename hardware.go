package hardware

import(
	"os"
	"../../def"
	"../elevator_io"
	"fmt"
	"../../ordermap"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}


func Hardware_init(new_floor chan int, new_order chan elevator_io.ButtonEvent){
	
	port := ":" + os.Args[2]
	myID := int((os.Args[1])[4])- 48 // ex.: os.Arg[1] = 49 
	fmt.Println("Elev: ",myID-1, "\n")
	def.MyID = myID

	ordermap.Current_order_floor = -1 

	elevator_io.Init(port, def.NUM_FLOORS)

	go elevator_io.Elevator_io_poll_floor_sensor(new_floor)
	go elevator_io.Elevator_io_poll_buttons(new_order)
	go Hardware_set_hall_lights() 

}


func Hardware_set_hall_lights(){ 
	for {
		for f := 0; f < def.NUM_FLOORS; f ++{
			for b := 0; b < def.NUM_BUTTONS-1; b ++ {
				if def.Num_elev_alive == 3 {
					if ordermap.Elev_info_matrix[0].Order_map[0][f][b] == 1 && ordermap.Elev_info_matrix[1].Order_map[1][f][b] == 1 && ordermap.Elev_info_matrix[2].Order_map[2][f][b] == 1 {
						mutex.Lock()
						elevator_io.SetButtonLamp(elevator_io.ButtonType(b), f, true)
						mutex.Unlock()
					}else {
						mutex.Lock()
						elevator_io.SetButtonLamp(elevator_io.ButtonType(b), f, false)
						mutex.Unlock()
					}
				}
				
				if def.Num_elev_alive == 2 {
					if ordermap.Elev_info_matrix[def.Elev_alive1-1].Order_map[def.Elev_alive1-1][f][b] == 1 && ordermap.Elev_info_matrix[def.Elev_alive2-1].Order_map[def.Elev_alive2-1][f][b] == 1 {
						mutex.Lock()
						elevator_io.SetButtonLamp(elevator_io.ButtonType(b), f, true)
						mutex.Unlock()
					}else { 
						mutex.Lock()
						elevator_io.SetButtonLamp(elevator_io.ButtonType(b), f, false)
						mutex.Unlock()
					}
				}					
			
				if def.Num_elev_alive == 1 && ordermap.Elev_info_matrix[def.Elev_alive-1].Order_map[def.Elev_alive-1][f][b] == 1 { 
					// Cannot guarantee the order to be completed. Therefore, the light is not turned on and rather being turned off.  
					mutex.Lock()
					elevator_io.SetButtonLamp(elevator_io.ButtonType(b), f, false)
					mutex.Unlock()
				}
			}
		}
		time.Sleep(500 * time.Millisecond)	
	}
}

