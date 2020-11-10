# go-dde

Reading data from the DDE server using GoLang

```
go get github.com/weibemoura/go-dde
```

#### Load methods "user32"
```
dde.Init()
```

#### Method callback for receive data
```
func ReceiveDataDde(item []byte, value []byte)  {
	fmt.Println(string(item), string(value))
}
```

#### Example
```
dde.Init()

client := dde.DdeClient{Callback: ReceiveDataDde}
defer client.Disconnect()

if client.Connect("profitchart", "cot") {
    //lastPrice := client.Request("WINFUT.ULT", 5000)
    //fmt.Printf("Last Price: %s\n", string(lastPrice))

    client.Advise("WINFUT.ULT", false)
    client.WinMSGLoop()
}
```