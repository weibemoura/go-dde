package dde

import (
	"log"
	"syscall"
	"unsafe"
)

var (
	DMLERR_NO_ERROR            = 0x0000
	DMLERR_ADVACKTIMEOUT       = 0x4000
	DMLERR_DATAACKTIMEOUT      = 0x4002
	DMLERR_DLL_NOT_INITIALIZED = 0x4003
	DMLERR_EXECACKTIMEOUT      = 0x4006
	DMLERR_NO_CONV_ESTABLISHED = 0x400A // (0x400A)
	DMLERR_POKEACKTIMEOUT      = 0x400B
	DMLERR_POSTMSG_FAILED      = 0x400C
	DMLERR_SERVER_DIED         = 0x400E

	CF_TEXT         = 1
	CF_BITMAP       = 2
	CF_METAFILEPICT = 3
	CF_SYLK         = 4
	CF_DIF          = 5
	CF_TIFF         = 6
	CF_OEMTEXT      = 7
	CF_DIB          = 8
	CF_PALETTE      = 9
	CF_PENDATA      = 10
	CF_RIFF         = 11
	CF_WAVE         = 12
	CF_UNICODETEXT  = 13
	CF_ENHMETAFILE  = 14
	CF_HDROP        = 15
	CF_LOCALE       = 16
	CF_DIBV5        = 17
	CF_MAX          = 18

	DDE_FACK          = 0x8000
	DDE_FBUSY         = 0x4000
	DDE_FDEFERUPD     = 0x4000
	DDE_FACKREQ       = 0x8000
	DDE_FRELEASE      = 0x2000
	DDE_FREQUESTED    = 0x1000
	DDE_FAPPSTATUS    = 0x00FF
	DDE_FNOTPROCESSED = 0x0000

	//DDE_FACKRESERVED = ~(DDE_FACK | DDE_FBUSY | DDE_FAPPSTATUS)
	//DDE_FADVRESERVED = ~(DDE_FACKREQ | DDE_FDEFERUPD)
	//DDE_FDATRESERVED = ~(DDE_FACKREQ | DDE_FRELEASE | DDE_FREQUESTED)
	//DDE_FPOKRESERVED = ~(DDE_FRELEASE)

	XTYPF_NOBLOCK = 0x0002
	XTYPF_NODATA  = 0x0004
	XTYPF_ACKREQ  = 0x0008

	XCLASS_MASK         = 0xFC00
	XCLASS_BOOL         = 0x1000
	XCLASS_DATA         = 0x2000
	XCLASS_FLAGS        = 0x4000
	XCLASS_NOTIFICATION = 0x8000

	XTYP_ERROR           = 0x0000 | XCLASS_NOTIFICATION | XTYPF_NOBLOCK
	XTYP_ADVDATA         = 0x0010 | XCLASS_FLAGS
	XTYP_ADVREQ          = 0x0020 | XCLASS_DATA | XTYPF_NOBLOCK
	XTYP_ADVSTART        = 0x0030 | XCLASS_BOOL
	XTYP_ADVSTOP         = 0x0040 | XCLASS_NOTIFICATION
	XTYP_EXECUTE         = 0x0050 | XCLASS_FLAGS
	XTYP_CONNECT         = 0x0060 | XCLASS_BOOL | XTYPF_NOBLOCK
	XTYP_CONNECT_CONFIRM = 0x0070 | XCLASS_NOTIFICATION | XTYPF_NOBLOCK
	XTYP_XACT_COMPLETE   = 0x0080 | XCLASS_NOTIFICATION
	XTYP_POKE            = 0x0090 | XCLASS_FLAGS
	XTYP_REGISTER        = 0x00A0 | XCLASS_NOTIFICATION | XTYPF_NOBLOCK
	XTYP_REQUEST         = 0x00B0 | XCLASS_DATA
	XTYP_DISCONNECT      = 0x00C0 | XCLASS_NOTIFICATION | XTYPF_NOBLOCK
	XTYP_UNREGISTER      = 0x00D0 | XCLASS_NOTIFICATION | XTYPF_NOBLOCK
	XTYP_WILDCONNECT     = 0x00E0 | XCLASS_DATA | XTYPF_NOBLOCK
	XTYP_MONITOR         = 0x00F0 | XCLASS_NOTIFICATION | XTYPF_NOBLOCK

	XTYP_MASK  = 0x00F0
	XTYP_SHIFT = 4

	TIMEOUT_ASYNC = 0xFFFFFFFF

	APPCLASS_STANDARD = 0x00000000
	APPCMD_CLIENTONLY = 0x00000010

	CP_WINANSI    = 1004
	CP_WINUNICODE = 1200
)

var (
	DdeAccessData          *syscall.Proc
	DdeClientTransaction   *syscall.Proc
	DdeConnect             *syscall.Proc
	DdeCreateDataHandle    *syscall.Proc
	DdeCreateStringHandleW *syscall.Proc
	DdeDisconnect          *syscall.Proc
	DdeGetLastError        *syscall.Proc
	DdeInitializeW         *syscall.Proc
	DdeFreeDataHandle      *syscall.Proc
	DdeFreeStringHandle    *syscall.Proc
	DdeQueryStringA        *syscall.Proc
	DdeUnaccessData        *syscall.Proc
	DdeUninitialize        *syscall.Proc

	GetMessageW      *syscall.Proc
	TranslateMessage *syscall.Proc
	DispatchMessageW *syscall.Proc
)

const (
	MAX_BUFFER_SIZE = 16
)

func Init() {
	user32 := syscall.MustLoadDLL("user32")
	DdeAccessData = user32.MustFindProc("DdeAccessData")
	DdeClientTransaction = user32.MustFindProc("DdeClientTransaction")
	DdeConnect = user32.MustFindProc("DdeConnect")
	DdeCreateDataHandle = user32.MustFindProc("DdeCreateDataHandle")
	DdeCreateStringHandleW = user32.MustFindProc("DdeCreateStringHandleW")
	DdeDisconnect = user32.MustFindProc("DdeDisconnect")
	DdeGetLastError = user32.MustFindProc("DdeGetLastError")
	DdeInitializeW = user32.MustFindProc("DdeInitializeW")
	DdeFreeDataHandle = user32.MustFindProc("DdeFreeDataHandle")
	DdeFreeStringHandle = user32.MustFindProc("DdeFreeStringHandle")
	DdeQueryStringA = user32.MustFindProc("DdeQueryStringA")
	DdeUnaccessData = user32.MustFindProc("DdeUnaccessData")
	DdeUninitialize = user32.MustFindProc("DdeUninitialize")

	GetMessageW = user32.MustFindProc("GetMessageW")
	TranslateMessage = user32.MustFindProc("TranslateMessage")
	DispatchMessageW = user32.MustFindProc("DispatchMessageW")
}

func (c *DdeClient) DdeCallback(wType int, wFmt int, hConv uintptr, hsz1 uintptr, hsz2 uintptr, hData uintptr, dwData1 int64, dwData2 int64) int {

	if wType == XTYP_ADVDATA {
		var pdwSize = 0
		pData, _, _ := DdeAccessData.Call(
			hData,
			uintptr(unsafe.Pointer(&pdwSize)))

		if int(pData) > 0 {
			value := *(*[MAX_BUFFER_SIZE]byte)(unsafe.Pointer(pData))

			var item [MAX_BUFFER_SIZE]byte
			dQuerySize, _, _ := DdeQueryStringA.Call(
				uintptr(c.IdInst),
				hsz2,
				uintptr(unsafe.Pointer(&item)),
				MAX_BUFFER_SIZE,
				uintptr(CP_WINANSI))

			c.Callback(
				item[:dQuerySize],
				value[:pdwSize-1])

			DdeUnaccessData.Call(hData)
			return DDE_FACK
		}
	}

	if wType == XTYP_DISCONNECT {
		log.Fatal("Disconnect notification received from server")
	}

	return 0
}

func (c *DdeClient) Connect(service string, topic string) bool {
	var strService, _ = syscall.UTF16PtrFromString(service)
	var strTopic, _ = syscall.UTF16PtrFromString(topic)

	res, _, err := DdeInitializeW.Call(
		uintptr(unsafe.Pointer(&c.IdInst)),
		syscall.NewCallback(c.DdeCallback),
		uintptr(APPCMD_CLIENTONLY),
		0)

	if int(res) != DMLERR_NO_ERROR {
		log.Fatal(err)
	}

	hszServName, _, _ := DdeCreateStringHandleW.Call(
		uintptr(c.IdInst),
		uintptr(unsafe.Pointer(strService)),
		uintptr(CP_WINUNICODE))

	hszTopic, _, _ := DdeCreateStringHandleW.Call(
		uintptr(c.IdInst),
		uintptr(unsafe.Pointer(strTopic)),
		uintptr(CP_WINUNICODE))

	var pConvContext interface{}
	hConv, _, _ := DdeConnect.Call(
		uintptr(c.IdInst),
		hszServName,
		hszTopic,
		uintptr(unsafe.Pointer(&pConvContext)))

	DdeFreeStringHandle.Call(uintptr(c.IdInst), hszTopic)
	DdeFreeStringHandle.Call(uintptr(c.IdInst), hszServName)

	if hConv > 0 {
		c.HConv = hConv
		return true
	}

	return false
}

func (c *DdeClient) Disconnect() {
	if int(c.HConv) > 0 {
		DdeDisconnect.Call(c.HConv)
	}

	if c.IdInst > 0 {
		DdeUninitialize.Call(uintptr(c.IdInst))
	}
}

func (c *DdeClient) Request(item string, timeout int) []byte {
	var strItem, _ = syscall.UTF16PtrFromString(item)

	hszItem, _, _ := DdeCreateStringHandleW.Call(
		uintptr(c.IdInst),
		uintptr(unsafe.Pointer(strItem)),
		uintptr(CP_WINUNICODE))

	var pdwResult int
	var lPByte interface{}
	hDdeData, _, _ := DdeClientTransaction.Call(
		uintptr(unsafe.Pointer(&lPByte)),
		0,
		c.HConv,
		hszItem,
		uintptr(CF_TEXT),
		uintptr(XTYP_REQUEST),
		uintptr(timeout),
		uintptr(unsafe.Pointer(&pdwResult)))

	DdeFreeStringHandle.Call(uintptr(c.IdInst), hszItem)

	if int(hDdeData) <= 0 {
		log.Fatalf("Unable to request item %d", int(hDdeData))
	}

	if timeout != TIMEOUT_ASYNC {
		var pdwSize = 0
		pData, _, _ := DdeAccessData.Call(
			hDdeData,
			uintptr(unsafe.Pointer(&pdwSize)))

		buffer := (*[MAX_BUFFER_SIZE]byte)(unsafe.Pointer(pData))[:pdwSize-1]
		if buffer == nil {
			DdeFreeDataHandle.Call(hDdeData)
			log.Fatalf("Unable to access data in request function %d", c.IdInst)
		}
		DdeUnaccessData.Call(hDdeData)
		return buffer
	}

	DdeFreeDataHandle.Call(hDdeData)

	return nil
}

func (c *DdeClient) Advise(item string, stop bool) {
	var strItem, _ = syscall.UTF16PtrFromString(item)

	var xtyp = XTYP_ADVSTART
	if stop {
		xtyp = XTYP_ADVSTOP
	}

	hszItem, _, _ := DdeCreateStringHandleW.Call(
		uintptr(c.IdInst),
		uintptr(unsafe.Pointer(strItem)),
		uintptr(CP_WINUNICODE))

	var pdwResult interface{}
	var lPByte interface{}
	hDdeData, _, _ := DdeClientTransaction.Call(
		uintptr(unsafe.Pointer(&lPByte)),
		0,
		c.HConv,
		hszItem,
		uintptr(CF_TEXT),
		uintptr(xtyp),
		uintptr(TIMEOUT_ASYNC),
		uintptr(unsafe.Pointer(&pdwResult)))

	DdeFreeStringHandle.Call(uintptr(c.IdInst), hszItem)

	if int(hDdeData) <= 0 {
		log.Fatalf("Unable to %s advise %d", "start", int(hDdeData))
	}

	DdeFreeDataHandle.Call(hDdeData)
}

func (c *DdeClient) Execute(command string) {
	var pData = []byte(command)
	var hSZ interface{}
	var pdwResult interface{}

	hDdeData, _, _ := DdeClientTransaction.Call(
		uintptr(unsafe.Pointer(&pData)),
		uintptr(len(pData) + 1),
		c.HConv,
		uintptr(unsafe.Pointer(&hSZ)),
		uintptr(CF_TEXT),
		uintptr(XTYP_EXECUTE),
		uintptr(TIMEOUT_ASYNC),
		uintptr(unsafe.Pointer(&pdwResult)))

	if int(hDdeData) <= 0 {
		log.Fatalf("Unable to send command %d", int(hDdeData))
	}

	DdeFreeDataHandle.Call(hDdeData)
}

func (c *DdeClient) Poke(item string, data []byte, timeout int) []byte {
	var strItem, _ = syscall.UTF16PtrFromString(item)

	hszItem, _, _ := DdeCreateStringHandleW.Call(
		uintptr(c.IdInst),
		uintptr(unsafe.Pointer(strItem)),
		uintptr(CP_WINUNICODE))

	var pdwResult = 0
	hDdeData, _, _ := DdeClientTransaction.Call(
		uintptr(unsafe.Pointer(&data)),
		uintptr(len(data) + 1),
		c.HConv,
		hszItem,
		uintptr(CF_TEXT),
		uintptr(XTYP_POKE),
		uintptr(timeout),
		uintptr(unsafe.Pointer(&pdwResult)))

	DdeFreeStringHandle.Call(uintptr(c.IdInst), hszItem)

	if int(hDdeData) <= 0 {
		log.Fatalf("Unable to poke to server %d", c.IdInst)
	}

	if timeout != TIMEOUT_ASYNC {
		var pdwSize = 0
		pData, _, _ := DdeAccessData.Call(
			hDdeData,
			uintptr(unsafe.Pointer(&pdwSize)))

		buffer := (*[MAX_BUFFER_SIZE]byte)(unsafe.Pointer(pData))[:pdwSize-1]
		if buffer == nil {
			DdeFreeDataHandle.Call(hDdeData)
			log.Fatalf("Unable to access data in poke function %d", c.IdInst)
		}
		DdeUnaccessData.Call(hDdeData)
		return buffer
	}

	DdeFreeDataHandle.Call(hDdeData)

	return nil

	return nil
}

func (c *DdeClient) WinMSGLoop() {
	var msg = Msg{}

	for {
		i, _, _ := GetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0, 0, 0)

		if int(i) <= 0 {
			break
		}

		TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}
