package main

import (
	"context"
	"fmt"
	"time"

	"github.com/finogeeks/ligase/sdk/go-ligase"
)

type Bot struct {
	devPrefix string
	userID    string
	token     string
	deviceID  string

	// sync
	Since string
}

func NewBot(devPrefix, userID string) *Bot {
	return &Bot{devPrefix: devPrefix, userID: userID}
}

func (b *Bot) GetUserID() string {
	return b.userID
}

func (b *Bot) GetToken() string {
	return b.token
}

func (b *Bot) GetDeviceID() string {
	return b.deviceID
}

func (b *Bot) Login() error {
	name := ""
	resp, err := ligase.Login(context.Background(), &ligase.PostLoginRequest{
		User:               b.userID,
		DeviceID:           b.devPrefix + b.userID,
		InitialDisplayName: &name,
	})
	if err != nil {
		return err
	}
	if resp.Header.StatusCode != 200 {
		if resp.Header.Error != nil {
			return resp.Header.Error
		} else {
			return fmt.Errorf("%s login failed, response statusCode: %d", b.userID, resp.Header.StatusCode)
		}
	}
	b.token = resp.AccessToken
	b.deviceID = resp.DeviceID
	return nil
}

func (b *Bot) Sync() (*ligase.SyncResponse, error) {
	var timeout time.Duration
	if b.Since != "" {
		timeout = time.Second * 15
	}
	resp, err := ligase.Sync(context.Background(), "", "", b.Since, "", "", timeout, b.token)
	if err != nil {
		return nil, err
	}
	if resp.Header.StatusCode != 200 {
		if resp.Header.Error != nil {
			return resp, resp.Header.Error
		} else {
			return resp, fmt.Errorf("%s sync failed, response statusCode: %d", b.userID, resp.Header.StatusCode)
		}
	}
	b.Since = resp.NextBatch

	return resp, nil
}

func (b *Bot) CreateRoom() (string, error) {
	resp, err := ligase.CreateRoom(context.Background(), &ligase.PostCreateRoomRequest{
		Name:       "",
		Topic:      `{"topic":"","group_property":{"who_can_invite":0,"who_can_talk":0,"enable_verification":false,"enable_e2ee":false,"enable_search":false,"version":1}}"`,
		Visibility: "public",
		Preset:     "public_chat",
		CreationContent: map[string]interface{}{
			"enable_favorite":  true,
			"enable_forward":   true,
			"enable_snapshot":  true,
			"enable_watermark": false,
			"is_direct":        false,
			"is_secret":        false,
			"m.federate":       false,
			"version":          1,
		},
	}, b.token)
	if err != nil {
		return "", err
	}
	if resp.Header.StatusCode != 200 {
		if resp.Header.Error != nil {
			return "", resp.Header.Error
		} else {
			return "", fmt.Errorf("%s createRoom failed, response statusCode: %d", b.userID, resp.Header.StatusCode)
		}
	}
	return resp.RoomID, nil
}

func (b *Bot) JoinRoom(roomID string) error {
	resp, err := ligase.JoinRoom(context.Background(), roomID, nil, b.token)
	if err != nil {
		return err
	}
	if resp.Header.StatusCode != 200 {
		if resp.Header.Error != nil {
			return resp.Header.Error
		} else {
			return fmt.Errorf("%s createRoom failed, response statusCode: %d", b.userID, resp.Header.StatusCode)
		}
	}
	return nil
}

func (b *Bot) SendMessage(roomID, body string) (string, error) {
	resp, err := ligase.SendMessageTxt(context.Background(), roomID, body, b.token)
	if err != nil {
		return "", err
	}
	if resp.Header.StatusCode != 200 {
		if resp.Header.Error != nil {
			return "", resp.Header.Error
		} else {
			return "", fmt.Errorf("%s createRoom failed, response statusCode: %d", b.userID, resp.Header.StatusCode)
		}
	}
	return resp.EventID, nil
}
