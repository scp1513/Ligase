package ligase

import (
	"encoding/json"

	"github.com/finogeeks/ligase/common/jsonerror"
	"github.com/finogeeks/ligase/model/syncapitypes"
	"github.com/finogeeks/ligase/plugins/message/external"
)

type Header struct {
	StatusCode int
	Error      *jsonerror.MatrixError
}

type (
	PostLoginRequest external.PostLoginRequest

	PostLoginResponse struct {
		external.PostLoginResponse
		Header Header `json:"-"`
	}

	SyncResponse struct {
		syncapitypes.Response
		Header Header `json:"-"`
	}

	PostCreateRoomRequest external.PostCreateRoomRequest

	PostCreateRoomResponse struct {
		external.PostCreateRoomResponse
		Header Header `json:"-"`
	}

	PostRoomsJoinByAliasRequest external.PostRoomsJoinByAliasRequest

	PostRoomsJoinByAliasResponse struct {
		external.PostRoomsJoinByAliasResponse
		Header Header `json:"-"`
	}

	PutRoomStateByTypeWithTxnIDResponse struct {
		external.PutRoomStateByTypeWithTxnIDResponse
		Header Header `json:"-"`
	}
)

func (r *PostLoginRequest) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *PostLoginResponse) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *SyncResponse) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *PostCreateRoomRequest) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *PostCreateRoomResponse) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *PostRoomsJoinByAliasRequest) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *PostRoomsJoinByAliasResponse) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func (r *PutRoomStateByTypeWithTxnIDResponse) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}
