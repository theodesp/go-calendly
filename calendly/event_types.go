package calendly

import (
	"context"
	"fmt"
	"strings"
)

const (
	eventTypesPath = "users/me/event_types"

	// Owner event type option
	IncludeTypeOwner IncludeType = "owner"
)

type EventTypesService ApiService

// Include Event type option
type IncludeType string


type EventTypesOpts struct {
	// request extra information about the entity that owns the Event Type,
	// by adding ?include=owner to the URL of the request.
	Include IncludeType `url:"include,omitempty"`
}

type eventTypesResponse struct {
	Data []*EventType `json:"data"`
}

type EventType struct {
	Type string `json:"type"`
	Id string `json:"id"`
	Attributes *EventTypeAttributes `json:"attributes"`
	Relationships *Relationships `json:"relationships,omitempty"`
}

type Relationships struct {
	Owner Owner `json:"owner"`
}

type Owner struct {
	Data Data `json:"data"`
}

type Data struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type EventTypeAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Duration    int64  `json:"duration"`
	Slug        string `json:"slug"`
	Color       string `json:"color"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	URL         string `json:"url"`
}

func (et *EventType) String() string  {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("EventType: id:%v attributes: ", et.Id))
	b.WriteString(fmt.Sprintf("Name:%v ", et.Attributes.Name))
	b.WriteString(fmt.Sprintf("Description:%v ", et.Attributes.Description))
	b.WriteString(fmt.Sprintf("Duration:%v ", et.Attributes.Duration))
	b.WriteString(fmt.Sprintf("Slug:%v ", et.Attributes.Slug))
	b.WriteString(fmt.Sprintf("Color:%v ", et.Attributes.Color))
	b.WriteString(fmt.Sprintf("Active:%v ", et.Attributes.Active))
	b.WriteString(fmt.Sprintf("CreatedAt:%v ", et.Attributes.CreatedAt))
	b.WriteString(fmt.Sprintf("UpdatedAt:%v ", et.Attributes.UpdatedAt))

	if et.Relationships != nil {
		b.WriteString(fmt.Sprintf("Owner Type:%v", et.Relationships.Owner.Data.Type))
		b.WriteString(fmt.Sprintf("Owner Id:%v", et.Relationships.Owner.Data.ID))
	}
	b.WriteString(fmt.Sprintf("URL:%v", et.Attributes.URL))

	return b.String()
}

// Event Types contain the most important configurations in Calendly.
// If you need some basic information about your event types, you can use this endpoint.
func (s *EventTypesService)List(ctx context.Context, opt *EventTypesOpts) ([]*EventType, *Response, error)  {
	u, err := addUrlOptions(eventTypesPath	, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.Get(u)
	if err != nil {
		return nil, nil, err
	}

	et := &eventTypesResponse{}
	resp, err := s.client.Do(ctx, req, et)
	if err != nil {
		return nil, resp, err
	}

	return et.Data, resp, nil
}