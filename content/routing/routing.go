// Copyright (C) 2020 Finogeeks Co., Ltd
//
// This program is free software: you can redistribute it and/or  modify
// it under the terms of the GNU Affero General Public License, version 3,
// as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package routing

import (
	"net/http"

	"github.com/finogeeks/ligase/common"
	"github.com/finogeeks/ligase/common/config"
	"github.com/finogeeks/ligase/common/uid"
	"github.com/finogeeks/ligase/content/download"
	"github.com/finogeeks/ligase/content/repos"
	"github.com/finogeeks/ligase/model/authtypes"
	"github.com/finogeeks/ligase/model/service"
	"github.com/finogeeks/ligase/model/types"
	"github.com/finogeeks/ligase/rpc"
	"github.com/finogeeks/ligase/skunkworks/log"
	mon "github.com/finogeeks/ligase/skunkworks/monitor/go-client/monitor"
	"github.com/gorilla/mux"
)

func Setup(
	apiMux *mux.Router,
	cfg *config.Dendrite,
	cacheIn service.Cache,
	feddomains *common.FedDomains,
	repo *repos.DownloadStateRepo,
	rpcClient rpc.RpcClient,
	consumer *download.DownloadConsumer,
	idg *uid.UidGenerator,
) {
	monitor := mon.GetInstance()
	histogram := monitor.NewLabeledHistogram(
		"http_requests_duration_milliseconds",
		[]string{"method", "path", "code"},
		[]float64{50.0, 100.0, 500.0, 1000.0, 2000.0, 5000.0},
	)

	prefixR0 := "/_matrix/media/r0"
	prefixV1 := "/_matrix/media/v1"

	muxR0 := apiMux.PathPrefix(prefixR0).Subrouter()
	muxV1 := apiMux.PathPrefix(prefixV1).Subrouter()

	processor := NewProcessor(cfg, histogram, repo, consumer, idg, []string{prefixR0, prefixV1})

	makeMediaAPI(muxR0, true, "/upload", processor.Upload, rpcClient, http.MethodPost, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/upload", processor.Upload, rpcClient, http.MethodPost, http.MethodOptions)

	makeMediaAPI(muxR0, false, "/download/{serverName}/{mediaId}", processor.Download, rpcClient, http.MethodGet, http.MethodOptions)
	makeMediaAPI(muxV1, false, "/download/{serverName}/{mediaId}", processor.Download, rpcClient, http.MethodGet, http.MethodOptions)

	makeMediaAPI(muxR0, false, "/thumbnail/{serverName}/{mediaId}", processor.Thumbnail, rpcClient, http.MethodGet, http.MethodOptions)
	makeMediaAPI(muxV1, false, "/thumbnail/{serverName}/{mediaId}", processor.Thumbnail, rpcClient, http.MethodGet, http.MethodOptions)

	makeMediaAPI(muxR0, true, "/favorite", processor.Favorite, rpcClient, http.MethodPost, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/favorite", processor.Favorite, rpcClient, http.MethodPost, http.MethodOptions)

	makeMediaAPI(muxR0, true, "/unfavorite/{netdiskID}", processor.Unfavorite, rpcClient, http.MethodDelete, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/unfavorite/{netdiskID}", processor.Unfavorite, rpcClient, http.MethodDelete, http.MethodOptions)

	makeMediaAPI(muxR0, true, "/forward/room/{roomID}", processor.SingleForward, rpcClient, http.MethodPost, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/forward/room/{roomID}", processor.SingleForward, rpcClient, http.MethodPost, http.MethodOptions)

	makeMediaAPI(muxR0, true, "/multi-forward", processor.MultiForward, rpcClient, http.MethodPost, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/multi-forward", processor.MultiForward, rpcClient, http.MethodPost, http.MethodOptions)

	makeMediaAPI(muxR0, true, "/multi-forward-res", processor.MultiResForward, rpcClient, http.MethodPost, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/multi-forward-res", processor.MultiResForward, rpcClient, http.MethodPost, http.MethodOptions)

	//emote
	//wait eif emote upload finish
	makeMediaAPI(muxR0, true, "/wait/emote", processor.WaitEmote, rpcClient, http.MethodGet, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/wait/emote", processor.WaitEmote, rpcClient, http.MethodGet, http.MethodOptions)
	//check emote is exsit
	makeMediaAPI(muxR0, true, "/check/emote/{serverName}/{mediaId}", processor.CheckEmote, rpcClient, http.MethodGet, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/check/emote/{serverName}/{mediaId}", processor.CheckEmote, rpcClient, http.MethodGet, http.MethodOptions)
	//favorite emote
	makeMediaAPI(muxR0, true, "/favorite/emote/{serverName}/{mediaId}", processor.FavoriteEmote, rpcClient, http.MethodPost, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/favorite/emote/{serverName}/{mediaId}", processor.FavoriteEmote, rpcClient, http.MethodPost, http.MethodOptions)
	//get emote list
	makeMediaAPI(muxR0, true, "/list/emote", processor.ListEmote, rpcClient, http.MethodGet, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/list/emote", processor.ListEmote, rpcClient, http.MethodGet, http.MethodOptions)
	//favorite file to emote
	makeMediaAPI(muxR0, true, "/favorite/fileemote/{serverName}/{mediaId}", processor.FavoriteFileEmote, rpcClient, http.MethodPost, http.MethodOptions)
	makeMediaAPI(muxV1, true, "/favorite/fileemote/{serverName}/{mediaId}", processor.FavoriteFileEmote, rpcClient, http.MethodPost, http.MethodOptions)

	fedV1 := apiMux.PathPrefix("/_matrix/federation/v1/media").Subrouter()
	makeFedAPI(fedV1, "/download/{serverName}/{mediaId}/{fileType}", processor.FedDownload, http.MethodGet, http.MethodOptions)
	makeFedAPI(fedV1, "/thumbnail/{serverName}/{mediaId}/{fileType}", processor.FedThumbnail, http.MethodGet, http.MethodOptions)
}

func verifyToken(rw http.ResponseWriter, req *http.Request, rpcCli rpc.RpcClient) (*authtypes.Device, bool) {
	token, err := common.ExtractAccessToken(req)
	if err != nil {
		log.Errorf("Content token error %s %v", req.RequestURI, err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error. " + err.Error()))
		return nil, false
	}

	rpcReq := types.VerifyTokenRequest{
		Token:      token,
		RequestURI: req.RequestURI,
	}
	content, err := rpcCli.VerifyToken(req.Context(), &rpcReq)
	if err != nil {
		log.Errorf("Content verify token error %s %v", req.RequestURI, err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error. " + err.Error()))
		return nil, false
	}
	if content.Error != "" {
		log.Errorf("Content verify token response error %s %v", req.RequestURI, err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(content.Error))
		return nil, false
	}

	device := content.Device

	return &device, true
}

func makeMediaAPI(r *mux.Router, atuh bool, url string, handler func(http.ResponseWriter, *http.Request, *authtypes.Device), rpcClient rpc.RpcClient, method ...string) {
	r.HandleFunc(url, func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				log.Errorf("Media API: %s panic %v", req.RequestURI, e)
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("Internal Server Error"))
			}
		}()
		if atuh {
			device, ok := verifyToken(rw, req, rpcClient)
			if !ok {
				return
			}
			handler(rw, req, device)
		} else {
			handler(rw, req, nil)
		}
	}).Methods(method...)
}

func makeFedAPI(r *mux.Router, url string, handler func(http.ResponseWriter, *http.Request), method ...string) {
	r.HandleFunc(url, func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				log.Errorf("Fed API: %s panic %v", req.RequestURI, e)
			}
		}()
		handler(rw, req)
	}).Methods(method...)
}
