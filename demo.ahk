#Requires AutoHotkey v2.0
#SingleInstance Force

^!q::ExitApp

OnExit ReleaseKeys

loopCount := 0

statusGui := Gui("+AlwaysOnTop +ToolWindow -MinimizeBox -MaximizeBox", "demo")
statusGui.MarginX := 12
statusGui.MarginY := 8
statusGui.SetFont("s10", "Segoe UI")
loopLabel := statusGui.AddText("w180", "Loop: 0")
stateLabel := statusGui.AddText("w180", "Starting in 10s")
exitButton := statusGui.AddButton("w180", "Exit")
exitButton.OnEvent("Click", (*) => ExitApp())
statusGui.OnEvent("Close", (*) => ExitApp())
statusGui.Show("x20 y20 NoActivate")
WinSetTransparent 230, "ahk_id " statusGui.Hwnd

Sleep 10000

Loop {
    loopCount += 1
    UpdateStatus("Holding W")

    Send "{w down}"
    Sleep 41000
    Send "{w up}"

    UpdateStatus("Waiting")
    Sleep 7000

    UpdateStatus("Pressing X")
    Send "x"
    Sleep 5000

    UpdateStatus("Pressing Enter")
    Send "{Enter}"

    UpdateStatus("Waiting 15s")
    Sleep 15000

    UpdateStatus("Pressing Enter")
    Send "{Enter}"
    Sleep 10000
}

UpdateStatus(state) {
    global loopCount, loopLabel, stateLabel
    loopLabel.Text := "Loop: " loopCount
    stateLabel.Text := state
}

ReleaseKeys(*) {
    Send "{w up}"
}
