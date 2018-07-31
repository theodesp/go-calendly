package calendly

import (
	"context"
	"net/url"
	"errors"
	"bytes"
	"fmt"
)

const (
	webhooksPath = "hooks"
	InviteeCreatedHookType EventHookType  = "invitee.created"
	InviteeCancelledHookType EventHookType  = "invitee.cancelled"
)

type WebhooksService apiService

type Webhook struct {
	ID string `json:"id"`
}

type EventHookType string

type WebhooksOpts struct {
	Url string
	Events []EventHookType
}


// Calendly supports webhooks which allow you to receive Calendly
// appointment data in real-time at a specified URL when a Calendly event is scheduled or cancelled.
//
// Specifically, you can subscribe to:
//
//  * Invitee Created Events (allowing you to receive notifications when a new Calendly event is created)
//  * Invitee Canceled Events (allowing you to receive notifications when a Calendly event is canceled)
// Creating a Webhook Subscription will not immediately trigger a webhook. So once it's set up, create or cancel an invitee to test it out.
func (s *WebhooksService)Create(ctx context.Context, opt *WebhooksOpts) (*Webhook, *Response, error)  {
	if opt == nil {
		return nil, nil, errors.New("go-calendly: webhooks.create required options")
	}

	u, err := url.Parse(opt.Url)
	if err != nil {
		return nil, nil, errors.New("go-calendly: webhooks.create url is not valid")
	}

	buf := bytes.NewBufferString("")
	buf.WriteString(url.QueryEscape(fmt.Sprintf("%v=%v", "url", u)))
	for _, ht := range opt.Events {
		buf.WriteString(url.QueryEscape(fmt.Sprintf("&%v=%v", "events[]", ht)))
	}

	req, err := s.client.Post(webhooksPath, buf.String())
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", textType)

	wh := &Webhook{}
	resp, err := s.client.Do(ctx, req, wh)
	if err != nil {
		return nil, resp, err
	}

	return wh, resp, nil
}
