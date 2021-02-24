package ligase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/finogeeks/ligase/common/jsonerror"
	"github.com/finogeeks/ligase/sdk/go-ligase/logger"
)

func Login(ctx context.Context, request *PostLoginRequest) (*PostLoginResponse, error) {
	url := ligaseURL + R0_PREFIX + "/login"
	payload, _ := json.Marshal(request)

	logger.GetLogger().Debugf("login req: %s, url: %s", string(payload), url)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	ts0 := time.Now()
	resp, err := httpCli.Do(req)
	ts1 := time.Now()
	logger.GetLogger().Tracef("login ligase spend time: %dms", ts1.Sub(ts0).Milliseconds())
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
		logger.GetLogger().Errorf("login ligase status code %d, body: %s", resp.StatusCode, data)
		e := jsonerror.MatrixError{}
		json.Unmarshal(data, &e)
		return &PostLoginResponse{Header: Header{StatusCode: resp.StatusCode, Error: &e}}, nil
	}

	if resp.Body == nil {
		return nil, errors.New("login ligase response body is nil, url = " + url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.GetLogger().Debugf("login ligase resp data: %s", string(data))

	var loginResp PostLoginResponse
	err = json.Unmarshal(data, &loginResp)
	if err != nil {
		return nil, err
	}
	logger.GetLogger().Debugf("login ligase resp: %s", &loginResp)
	loginResp.Header.StatusCode = resp.StatusCode

	return &loginResp, nil
}
