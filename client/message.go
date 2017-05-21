// Code generated by goagen v1.2.0, DO NOT EDIT.
//
// API "Chat API": message Resource Client
//
// Command:
// $ goagen
// --design=github.com/m0a-mystudy/goa-chat/design
// --out=$(GOPATH)/src/github.com/m0a-mystudy/goa-chat
// --version=v1.2.0-dirty

package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ListMessagePath computes a request path to the list action of message.
func ListMessagePath(roomID int) string {
	param0 := strconv.Itoa(roomID)

	return fmt.Sprintf("/api/rooms/%s/messages", param0)
}

// Retrieve all messages.
func (c *Client) ListMessage(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewListMessageRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListMessageRequest create the request corresponding to the list action endpoint of the message resource.
func (c *Client) NewListMessageRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// PostMessagePath computes a request path to the post action of message.
func PostMessagePath(roomID int) string {
	param0 := strconv.Itoa(roomID)

	return fmt.Sprintf("/api/rooms/%s/messages", param0)
}

// Create new message
func (c *Client) PostMessage(ctx context.Context, path string, payload *MessagePayload, contentType string) (*http.Response, error) {
	req, err := c.NewPostMessageRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewPostMessageRequest create the request corresponding to the post action endpoint of the message resource.
func (c *Client) NewPostMessageRequest(ctx context.Context, path string, payload *MessagePayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	if c.BasicAuthSigner != nil {
		c.BasicAuthSigner.Sign(req)
	}
	return req, nil
}

// ShowMessagePath computes a request path to the show action of message.
func ShowMessagePath(roomID int, messageID int) string {
	param0 := strconv.Itoa(roomID)
	param1 := strconv.Itoa(messageID)

	return fmt.Sprintf("/api/rooms/%s/messages/%s", param0, param1)
}

// Retrieve message with given id
func (c *Client) ShowMessage(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewShowMessageRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewShowMessageRequest create the request corresponding to the show action endpoint of the message resource.
func (c *Client) NewShowMessageRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
