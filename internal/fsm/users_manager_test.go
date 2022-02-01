package fsm

import (
	"testing"
)

const (
	testStateIdle State = iota
	testStateWaiting

	testEventStartWait Event = iota
	testEventStopWait
	testEventCancel
)

var (
	testFsm = NewFSM(Transitions{
		testStateIdle: {
			testEventStartWait: testStateWaiting,
		},
		testStateWaiting: {
			testEventStopWait: testStateIdle,
			testEventCancel:   testStateIdle,
		},
	})
)

func TestUsersManager_GetState(t *testing.T) {
	type fields struct {
		idleState   State
		Fsm         FSM
		usersStates map[int64]State
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantState State
	}{
		{
			name: "No user, return idle state",
			fields: fields{
				idleState:   testStateIdle,
				Fsm:         testFsm,
				usersStates: map[int64]State{},
			},
			args:      args{userId: 1},
			wantState: testStateIdle,
		},
		{
			name: "Return user state",
			fields: fields{
				idleState: testStateIdle,
				Fsm:       testFsm,
				usersStates: map[int64]State{
					1: testStateWaiting,
				},
			},
			args:      args{userId: 1},
			wantState: testStateWaiting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := &UsersManager{
				idleState:   tt.fields.idleState,
				Fsm:         tt.fields.Fsm,
				usersStates: tt.fields.usersStates,
			}
			if gotState := um.GetState(tt.args.userId); gotState != tt.wantState {
				t.Errorf("GetState() = %v, want %v", gotState, tt.wantState)
			}
		})
	}
}

func TestUsersManager_SendEventE(t *testing.T) {
	type fields struct {
		idleState   State
		Fsm         FSM
		usersStates map[int64]State
	}
	type args struct {
		userId int64
		event  Event
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantNewState State
		wantErr      error
	}{
		{
			name: "No user",
			fields: fields{
				idleState:   testStateIdle,
				Fsm:         testFsm,
				usersStates: map[int64]State{},
			},
			args:         args{userId: 1, event: testEventStartWait},
			wantNewState: testStateWaiting,
		},
		{
			name: "Transition by event",
			fields: fields{
				idleState: testStateIdle,
				Fsm:       testFsm,
				usersStates: map[int64]State{
					1: testStateIdle,
				},
			},
			args:         args{userId: 1, event: testEventStartWait},
			wantNewState: testStateWaiting,
		},
		{
			name: "Keep state on invalid event",
			fields: fields{
				idleState: testStateIdle,
				Fsm:       testFsm,
				usersStates: map[int64]State{
					1: testStateIdle,
				},
			},
			args:         args{userId: 1, event: testEventCancel},
			wantNewState: testStateIdle,
			wantErr:      ErrEventForCurrentStateNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			um := &UsersManager{
				idleState:   tt.fields.idleState,
				Fsm:         tt.fields.Fsm,
				usersStates: tt.fields.usersStates,
			}
			gotNewState, err := um.SendEventE(tt.args.userId, tt.args.event)
			if err != tt.wantErr {
				t.Errorf("SendEventE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNewState != tt.wantNewState {
				t.Errorf("SendEventE() gotNewState = %v, want %v", gotNewState, tt.wantNewState)
			}
		})
	}
}
