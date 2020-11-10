package main

import (
	"fmt"
	"github.com/weibemoura/go-dde/dde"
)

func ReceiveDataDde(item []byte, value []byte)  {
	fmt.Println(string(item), string(value))
}

func main()  {
	dde.Init()

	client := dde.DdeClient{Callback: ReceiveDataDde}
	defer client.Disconnect()

	if client.Connect("profitchart", "cot") {
		//lastPrice := client.Request("WINFUT.ULT", 5000)
		//fmt.Printf("Last Price: %s\n", string(lastPrice))

		client.Advise("WINFUT.ULT", false)
		client.WinMSGLoop()
	}
}
