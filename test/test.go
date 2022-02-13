package main

import (
	"fmt"
	"os"
)

func main(){
	f, err := os.OpenFile("IU(아이유) 'Blueming(블루밍)' 라이브🎤🎤(밴드ver.)  가사  스페셜클립  Special Clip  LYRICS [4K].webm", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err:",err)
	}
	defer f.Close()
}
