package event

import (
	"onTuneKubeManager/common"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

const (
	MANAGER_SIGNAL   = "OnTuneKubeManagerSignal"
	TERMINATE_SIGNAL = "OnTuneKubeManagerSignalTerminate"
)

// SetEvent is a function that defines setEvent and terminateEvent.
func SetEvent() {
	var hevent windows.Handle
	hevent, err := windows.OpenEvent(windows.EVENT_ALL_ACCESS, false, windows.StringToUTF16Ptr(MANAGER_SIGNAL))
	if err != nil {
		hevent, err = windows.CreateEvent(nil, 0, 0, windows.StringToUTF16Ptr(MANAGER_SIGNAL))
		if err != nil {
			return
		}
	}

	var htevent windows.Handle
	htevent, err = windows.OpenEvent(windows.EVENT_ALL_ACCESS, false, windows.StringToUTF16Ptr(TERMINATE_SIGNAL))
	if err != nil {
		htevent, err = windows.CreateEvent(nil, 0, 0, windows.StringToUTF16Ptr(TERMINATE_SIGNAL))
		if err != nil {
			return
		}
	}

	common.LogManager.WriteLog("Start onTune Kube Manager by Windows Event")
	go SetWindowsEvent(hevent)
	go TerminateWindowsEvent(htevent)
}

// SetWindowsEvent is a function that sets the event in Windows.
// hevent is a handle of event.
func SetWindowsEvent(hevent windows.Handle) {
	for {
		err := windows.SetEvent(hevent)
		if err != nil {
			return
		}
		time.Sleep(time.Second * time.Duration(1))
	}
}

// TerminateWindowsEvent is a function that terminates the event in Windows.
// htevent is a handle of terminated event.
func TerminateWindowsEvent(htevent windows.Handle) {
	_, err := windows.WaitForSingleObject(htevent, windows.INFINITE)
	if err != nil {
		common.LogManager.WriteLog(err.Error())
	}
	common.LogManager.WriteLog("Terminate onTune Kube Manager by Windows Event")
	os.Exit(0)
}
