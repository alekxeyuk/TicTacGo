// Code generated by "stringer -type=RoomState,PlayerSign"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EMPTY_ROOM-0]
	_ = x[GAME_IN_PROGRESS-1]
	_ = x[GAME_FINISHED-2]
}

const _RoomState_name = "EMPTY_ROOMGAME_IN_PROGRESSGAME_FINISHED"

var _RoomState_index = [...]uint8{0, 10, 26, 39}

func (i RoomState) String() string {
	if i >= RoomState(len(_RoomState_index)-1) {
		return "RoomState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RoomState_name[_RoomState_index[i]:_RoomState_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EMPTY_CELL-0]
	_ = x[PLAYER_X-1]
	_ = x[PLAYER_O-2]
}

const _PlayerSign_name = "EMPTY_CELLPLAYER_XPLAYER_O"

var _PlayerSign_index = [...]uint8{0, 10, 18, 26}

func (i PlayerSign) String() string {
	if i >= PlayerSign(len(_PlayerSign_index)-1) {
		return "PlayerSign(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PlayerSign_name[_PlayerSign_index[i]:_PlayerSign_index[i+1]]
}
