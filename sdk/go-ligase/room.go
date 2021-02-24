package ligase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/finogeeks/ligase/common/jsonerror"
	"github.com/finogeeks/ligase/sdk/go-ligase/logger"
)

func CreateRoom(ctx context.Context, request *PostCreateRoomRequest, token string) (*PostCreateRoomResponse, error) {
	url := fmt.Sprintf(ligaseURL+R0_PREFIX+"/createRoom?access_token=%s", token)
	payload, _ := json.Marshal(request)
	logger.GetLogger().Debugf("createRoom request: %s url: %s", string(payload), url)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	ts0 := time.Now()
	resp, err := httpCli.Do(req)
	ts1 := time.Now()
	logger.GetLogger().Tracef("createRoom spend time: %dms", ts1.Sub(ts0).Milliseconds())
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		var data []byte
		if resp.Body != nil {
			data, _ = ioutil.ReadAll(resp.Body)
		}
		logger.GetLogger().Errorf("createRoom status code %d, body: %s", resp.StatusCode, data)
		e := jsonerror.MatrixError{}
		json.Unmarshal(data, &e)
		return &PostCreateRoomResponse{Header: Header{StatusCode: resp.StatusCode, Error: &e}}, nil
	}

	if resp.Body == nil {
		return &PostCreateRoomResponse{Header: Header{StatusCode: resp.StatusCode}}, errors.New("createRoom response body is nil, url = " + url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &PostCreateRoomResponse{Header: Header{StatusCode: resp.StatusCode}}, err
	}

	logger.GetLogger().Debugf("createRoom resp data: %s", string(data))

	var createRoomResp PostCreateRoomResponse
	err = json.Unmarshal(data, &createRoomResp)
	if err != nil {
		return &PostCreateRoomResponse{Header: Header{StatusCode: resp.StatusCode}}, err
	}
	logger.GetLogger().Debugf("createRoom resp: %s", &createRoomResp)
	createRoomResp.Header.StatusCode = resp.StatusCode

	return &createRoomResp, nil
}

func JoinRoom(ctx context.Context, roomIDOrAlias string, content map[string]interface{}, token string) (*PostRoomsJoinByAliasResponse, error) {
	url := fmt.Sprintf(ligaseURL+R0_PREFIX+"/join/%s?access_token=%s", roomIDOrAlias, token)
	var payload []byte
	if content != nil {
		payload, _ = json.Marshal(content)
	} else {
		payload = []byte("{}")
	}
	logger.GetLogger().Debugf("joinRoom content: %s url: %s", string(payload), url)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	ts0 := time.Now()
	resp, err := httpCli.Do(req)
	ts1 := time.Now()
	logger.GetLogger().Tracef("joinRoom spend time: %dms", ts1.Sub(ts0).Milliseconds())
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		var data []byte
		if resp.Body != nil {
			data, _ = ioutil.ReadAll(resp.Body)
		}
		logger.GetLogger().Errorf("joinRoom status code %d, body: %s", resp.StatusCode, data)
		e := jsonerror.MatrixError{}
		json.Unmarshal(data, &e)
		return &PostRoomsJoinByAliasResponse{Header: Header{StatusCode: resp.StatusCode, Error: &e}}, nil
	}

	if resp.Body == nil {
		return &PostRoomsJoinByAliasResponse{Header: Header{StatusCode: resp.StatusCode}}, errors.New("joinRoom response body is nil, url = " + url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &PostRoomsJoinByAliasResponse{Header: Header{StatusCode: resp.StatusCode}}, err
	}

	logger.GetLogger().Debugf("joinRoom resp data: %s", string(data))

	var joinRoomResp PostRoomsJoinByAliasResponse
	err = json.Unmarshal(data, &joinRoomResp)
	if err != nil {
		return &PostRoomsJoinByAliasResponse{Header: Header{StatusCode: resp.StatusCode}}, err
	}
	logger.GetLogger().Debugf("joinRoom resp: %s", &joinRoomResp)
	joinRoomResp.Header.StatusCode = resp.StatusCode

	return &joinRoomResp, nil
}

var (
	messageCounter int64
)

func SendMessage(ctx context.Context, roomID string, content map[string]interface{}, token string) (*PutRoomStateByTypeWithTxnIDResponse, error) {
	url := fmt.Sprintf(ligaseURL+R0_PREFIX+"/rooms/%s/send/m.room.message/m%d.%d?access_token=%s", roomID, time.Now().UnixNano()/1000000, atomic.AddInt64(&messageCounter, 1), token)
	payload, _ := json.Marshal(content)
	logger.GetLogger().Debugf("sendMessage content: %s url: %s", string(payload), url)

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	ts0 := time.Now()
	resp, err := httpCli.Do(req)
	ts1 := time.Now()
	logger.GetLogger().Tracef("sendMessage spend time: %dms", ts1.Sub(ts0).Milliseconds())
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		var data []byte
		if resp.Body != nil {
			data, _ = ioutil.ReadAll(resp.Body)
		}
		logger.GetLogger().Errorf("sendMessage status code %d, body: %s", resp.StatusCode, data)
		e := jsonerror.MatrixError{}
		json.Unmarshal(data, &e)
		return &PutRoomStateByTypeWithTxnIDResponse{Header: Header{StatusCode: resp.StatusCode, Error: &e}}, nil
	}

	if resp.Body == nil {
		return &PutRoomStateByTypeWithTxnIDResponse{Header: Header{StatusCode: resp.StatusCode}}, errors.New("sendMessage response body is nil, url = " + url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &PutRoomStateByTypeWithTxnIDResponse{Header: Header{StatusCode: resp.StatusCode}}, err
	}

	logger.GetLogger().Debugf("sendMessage resp data: %s", string(data))

	var sendMessageResp PutRoomStateByTypeWithTxnIDResponse
	err = json.Unmarshal(data, &sendMessageResp)
	if err != nil {
		return &PutRoomStateByTypeWithTxnIDResponse{Header: Header{StatusCode: resp.StatusCode}}, err
	}
	logger.GetLogger().Debugf("sendMessage resp: %s", &sendMessageResp)
	sendMessageResp.Header.StatusCode = resp.StatusCode

	return &sendMessageResp, nil
}

func SendMessageTxt(ctx context.Context, roomID, body, token string) (*PutRoomStateByTypeWithTxnIDResponse, error) {
	return SendMessage(ctx, roomID, map[string]interface{}{
		"msgtype": "m.text",
		"body":    body,
	}, token)
}
