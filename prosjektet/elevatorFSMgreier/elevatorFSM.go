//Data:
//Pickup request array fra prAssigner - [[up,down],[up,down],[up,down],[up,down]]
//Destination request array fra MODULUTENNAVN - [bool,bool,bool,bool]
//Retning heisen kjører i - direction [up,down]

//---------------------------------------------

//while løkke som går igjennom alt forbi dette punktet --------

//FSM

//Om heisen starter opp mellom etasjer
//fsm_onInitBetweenFloors(void)
//endre motorretning til nedover
//Sette heisretning til nedover
//Sette heisbehavior til moving

//Når heisen får et buttonpress
//fsm_onRequestButtonPress(int btn_floor, Button btn_type)   ER VOID
//switch(elevator.behavior)
//case EB_DoorOpen:
	//if(requests_shouldClearImmediately(elevator, btn_floor, btn_type)){
	//timer_start(elevator.config.doorOpenDuration_s);
    //} else {
        //PR for etasjen det trykkes i i retningen det trykkes i skal settes til true
    //}
        //break;
//case EB_Moving:
	s



void fsm_onRequestButtonPress(int btn_floor, Button btn_type){
    printf("\n\n%s(%d, %s)\n", __FUNCTION__, btn_floor, elevio_button_toString(btn_type));
    elevator_print(elevator);
    
    switch(elevator.behaviour){
    case EB_DoorOpen:
        if(requests_shouldClearImmediately(elevator, btn_floor, btn_type)){
            timer_start(elevator.config.doorOpenDuration_s);
        } else {
            elevator.requests[btn_floor][btn_type] = 1;
        }
        break;

    case EB_Moving:
        elevator.requests[btn_floor][btn_type] = 1;
        break;
        
    case EB_Idle:    
        elevator.requests[btn_floor][btn_type] = 1;
        DirnBehaviourPair pair = requests_chooseDirection(elevator);
        elevator.dirn = pair.dirn;
        elevator.behaviour = pair.behaviour;
        switch(pair.behaviour){
        case EB_DoorOpen:
            outputDevice.doorLight(1);
            timer_start(elevator.config.doorOpenDuration_s);
            elevator = requests_clearAtCurrentFloor(elevator);
            break;

        case EB_Moving:
            outputDevice.motorDirection(elevator.dirn);
            break;
            
        case EB_Idle:
            break;
        }
        break;
    }
    
    setAllLights(elevator);
    
    printf("\nNew state:\n");
    elevator_print(elevator);
}


//Hva heisen skal gjøre når den ankommer en etasje
//fsm_onFloorArrival(int newFloor)  newFloor er etasjen vi har kommet til
// sette heisens etasjevariabel lik newFloor
// switch med moving og default som states:
//     if moving:
//          stopp motor, skru på open door light, fjerne DR på etasjen vi er i og PR i retningen vi går i,
//          starte timer for åpen dør, sette elevator behavior til doorOpen
//          break;
//     if default: break;


//If elevator.behavior = idle
//      check if floor

//If elevator.behavior = moving
//      hente etasjevariabel imens den kjører
//      vite retningen heisen kjører i
//      hvis etasjevariabel er på listen over etasjer som skal stoppes i, stopp

//If elevator.behavior = door open
// vente til dør lukkes, så sjekke om heis er idle eller moving?