package def

const (
	NUM_FLOORS          = 4
	NUM_BUTTONS         = 3
	NUM_ELEV            = 3
	TIME_DOOR_OPEN      = 3
	TIMEOUT_SERVE_ORDER =10
	EMERGENCY_STOP_TIME = 3
)

var MyID           int    // Elevator 1, 2, 3
var Elev_alive     int    // The ID to an elevator that is guaranteed to be alive 
var Num_elev_alive int 
var Elev_alive1    int	  // For control of which elevators are alive if there is only two elevetaors alive
var Elev_alive2    int    // For control of which elevators are alive if there is only two elevetaors alive