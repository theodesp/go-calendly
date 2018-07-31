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
	getWebhookpath = "hooks/%v"
)

type WebhooksService apiService

type Webhook struct {
	Type       string     `json:"type"`
	ID         int64      `json:"id"`
	Attributes *WebhookAttributes `json:"attributes,omitempty"`
}

type WebhookAttributes struct {
	URL       string  `json:"url"`
	CreatedAt string  `json:"created_at"`
	State     string  `json:"state"`
	Events    []EventHookType `json:"events"`
}

type EventHookType string
const (
	InviteeCreatedHookType EventHookType  = "invitee.created"
	InviteeCancelledHookType EventHookType  = "invitee.cancelled"
)

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

// You can view your current Webhook Subscriptions.
//
// Using this endpoint will list up to the first 100 Webhook Subscriptions,
// and it will order active subscriptions first.
func (s *WebhooksService)List(ctx context.Context) ([]*Webhook, *Response, error)  {
	req, err := s.client.Get(webhooksPath)
	if err != nil {
		return nil, nil, err
	}

	wh := &webhookListResponse{}
	resp, err := s.client.Do(ctx, req, wh)
	if err != nil {
		return nil, resp, err
	}

	return wh.Webhooks, resp, nil
}

// Any of your Webhook Subscriptions can be accessed by ID using this endpoint.
func (s *WebhooksService)GetByID(ctx context.Context, id int64) (*Webhook, *Response, error)  {
	req, err := s.client.Get(fmt.Sprintf(getWebhookpath, id))
	if err != nil {
		return nil, nil, err
	}

	wh := &webhookResponse{}
	resp, err := s.client.Do(ctx, req, wh)
	if err != nil {
		return nil, resp, err
	}

	return wh.Webhook, resp, nil
}

type webhookResponse struct {
	Webhook *Webhook `json:"data"`
}

type webhookListResponse struct {
	Webhooks []*Webhook `json:"data"`
}
