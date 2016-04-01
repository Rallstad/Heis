package orders

import ( 
	. "fmt"
	. "../driver"
)



const N_ELEV = 3



var orders = [N_FLOOR][N_BUTTON_TYPES]int{
	{0,-1,0},
	{0,0,0},
	{0,0,0},
	{-1,0,0},
}

func Register_order(){
	
	for i := 0; i < N_FLOOR; i++{
		for j:=0;j<N_BUTTON_TYPES;j++{
			if Elev_get_button_signal(j,i) > 0{
				orders[i][j]=1
			}
		}
	}
	
}

func Set_light(){
	for i := 0; i < N_FLOOR; i++{
		for j:=0;j<N_BUTTON_TYPES;j++{
			if orders[i][j] >0{
				Elev_set_button_lamp(j,i,1)
			}
		}
	}
}

func Print_ext_orders(){
	Println("Current orders")
	for i := 0; i < N_FLOOR; i++{
		for j:=0;j<N_BUTTON_TYPES;j++{
			Print(orders[i][j])

		}
		Println("\n")
	}
	Println("\n")
	
}



