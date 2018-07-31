package calendly

import (
	"context"
	"bytes"
	"fmt"
)

const (
	aboutMePath = "users/me"
)

type UsersService apiService

type AboutMeResponse struct {
	AboutMe *AboutMe `json:"data"`
}

type AboutMe struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Attributes *UserAttributes `json:"attributes,omitempty"`
}

type UserAttributes struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Email     string `json:"email"`
	URL       string `json:"url"`
	Timezone  string `json:"timezone"`
	Avatar    *Avatar `json:"avatar,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Avatar struct {
	URL string `json:"url"`
}

func (a *AboutMe) String() string  {
	b := bytes.NewBufferString("")
	b.WriteString(fmt.Sprintf("About Me: id:%v attributes: ", a.ID))
	b.WriteString(fmt.Sprintf("Name:%v ", a.Attributes.Name))
	b.WriteString(fmt.Sprintf("Email:%v ", a.Attributes.Email))
	b.WriteString(fmt.Sprintf("Slug:%v ", a.Attributes.Slug))
	b.WriteString(fmt.Sprintf("Timezone:%v ", a.Attributes.Timezone))
	b.WriteString(fmt.Sprintf("CreatedAt:%v ", a.Attributes.CreatedAt))
	b.WriteString(fmt.Sprintf("UpdatedAt:%v ", a.Attributes.UpdatedAt))

	if a.Attributes.Avatar != nil {
		b.WriteString(fmt.Sprintf("Avatar:%v ", a.Attributes.Avatar.URL))
	}
	b.WriteString(fmt.Sprintf("URL:%v ", a.Attributes.URL))

	return b.String()
}

// Use this endpoint to request basic information about yourself.
// This might be helpful if you're building functionality for multiple Calendly users.
func (s *UsersService)AboutMe(ctx context.Context) (*AboutMe, *Response, error)  {
	req, err := s.client.Get(aboutMePath)
	if err != nil {
		return nil, nil, err
	}

	a := &AboutMeResponse{}
	resp, err := s.client.Do(ctx, req, a)
	if err != nil {
		return nil, resp, err
	}

	return a.AboutMe, resp, nil
}
