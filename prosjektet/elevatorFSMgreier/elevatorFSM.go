//Data:
//Pickup request array fra prAssigner - [[up,down],[up,down],[up,down],[up,down]]
//Destination request array fra MODULUTENNAVN - [bool,bool,bool,bool]
//Retning heisen kjører i - direction [up,down]

//---------------------------------------------

//while løkke som går igjennom alt forbi dette punktet --------

//FSM

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