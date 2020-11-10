package dde

type Point struct {
	x int
	y int
}

type Msg struct {
	hWnd    int
	message int
	wParam  int
	lParam  int
	time    int
	pt      Point
}

type DdeClient struct {
	IdInst int
	HConv  uintptr
	Callback func(item []byte, value []byte)
}