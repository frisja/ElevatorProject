package timer

import (
	"time"
	"../ordermap"
	"../def"
)

func Timer_is_elev_completing_order(get_current_floor chan int, set_dir chan int, peerTxEnable chan bool) {
 	
 	timer1 := time.NewTimer(def.TIMEOUT_SERVE_ORDER * time.Second)
 	
 	for {
 		select{
 		case <-timer1.C:
			set_dir <- ordermap.Current_dir
			peerTxEnable <- false
 			return 
 		
 		case newfloor :=<- get_current_floor:
 			if ordermap.Current_order_floor == newfloor {
 				return 
 			}
 		}
 		time.Sleep(20 * time.Millisecond) 
 	}		
}


func Timer_get_current_floor(get_current_floor chan int) {
	for {
		get_current_floor <- ordermap.Current_floor
	}
}
