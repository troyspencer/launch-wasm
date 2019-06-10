package world

func (worldState *WorldState) TriggerUpdateLaunches() {
	launched := worldState.Doc.Call("createEvent", "Event")
	launched.Call("initEvent", "updateLaunches")
	launched.Set("launches", worldState.Launches)
	worldState.Doc.Call("dispatchEvent", launched)
}

func (worldState *WorldState) TriggerResetTimer() {
	resetTimer := worldState.Doc.Call("createEvent", "Event")
	resetTimer.Call("initEvent", "resetTimer")
	worldState.Doc.Call("dispatchEvent", resetTimer)
}
