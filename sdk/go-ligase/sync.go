package ligase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/finogeeks/ligase/sdk/go-ligase/logger"
)

func Sync(ctx context.Context, filter, from, since, setPresence, fullstate string, timeout time.Duration, token string) (*SyncResponse, error) {
	url := fmt.Sprintf(ligaseURL+R0_PREFIX+"/sync?filter=%s&from=%s&since=%sset_presence=%s&full_state=%s&timeout=%d&access_token=%s",
		filter, from, since, setPresence, fullstate, timeout.Milliseconds(), token)
	logger.GetLogger().Debugf("sync url: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	ts0 := time.Now()
	resp, err := httpCli.Do(req)
	ts1 := time.Now()
	logger.GetLogger().Tracef("sync spend time: %dms", ts1.Sub(ts0).Milliseconds())
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
		logger.GetLogger().Errorf("sync status code %d, body: %s", resp.StatusCode, data)
		var syncResp SyncResponse
		json.Unmarshal(data, &syncResp)
		logger.GetLogger().Debugf("sync resp: %s", &syncResp)
		syncResp.Header.StatusCode = resp.StatusCode
		return &syncResp, nil
	}

	if resp.Body == nil {
		return &SyncResponse{Header: Header{StatusCode: resp.StatusCode}}, errors.New("sync response body is nil, url = " + url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &SyncResponse{Header: Header{StatusCode: resp.StatusCode}}, err
	}

	logger.GetLogger().Debugf("sync resp data: %s", string(data))

	var syncResp SyncResponse
	err = json.Unmarshal(data, &syncResp)
	if err != nil {
		return &SyncResponse{Header: Header{StatusCode: resp.StatusCode}}, err
	}
	logger.GetLogger().Debugf("sync resp: %s", &syncResp)
	syncResp.Header.StatusCode = resp.StatusCode

	return &syncResp, nil
}
