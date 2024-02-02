package todoistclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CommentCount   int    `json:"comment_count"`
	Color          string `json:"color"`
	IsShared       bool   `json:"is_shared"`
	Order          int    `json:"order"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	ViewStyle      string `json:"view_style"`
	URL            string `json:"url"`
	ParentID       string `json:"parent_id"`
}

func (c *Client) GetProject(ctx context.Context, projectId string) (*Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s", c.BaseURL, projectId), nil)
	log.WithFields(log.Fields{
		"projectId": projectId,
	}).Info("Reading project")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req = req.WithContext(ctx)
	res := Project{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"Project": fmt.Sprintf("%+v", res),
	}).Debug("Project read")
	return &res, nil
}

type CreateProject struct {
	// Required fields
	Name *string `json:"name"`
	// Optional fields
	ParentID   *string `json:"parent_id,omitempty"`
	Color      *string `json:"color,omitempty"`
	IsFavorite *bool   `json:"is_favorite,omitempty"`
	ViewStyle  *string `json:"view_style,omitempty"`
}

func (c *Client) CreateProject(ctx context.Context, createProject CreateProject) (*Project, error) {
	payload, err := json.Marshal(createProject)
	log.WithFields(log.Fields{
		"payload": string(payload),
	}).Info("Creating project")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects", c.BaseURL), bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req = req.WithContext(ctx)
	res := Project{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Project": fmt.Sprintf("%+v", res),
	}).Debug("Project created")
	return &res, nil
}

type UpdateProject struct {
	ID         *string
	Name       *string `json:"name,omitempty"`
	Color      *string `json:"color,omitempty"`
	IsFavorite *bool   `json:"is_favorite,omitempty"`
	ViewStyle  *string `json:"view_style,omitempty"`
}

func (c *Client) UpdateProject(ctx context.Context, updateProject UpdateProject) (*Project, error) {
	payload, err := json.Marshal(updateProject)
	if updateProject.ID == nil {
		return nil, fmt.Errorf("missing project id")
	}
	log.WithFields(log.Fields{
		"payload":   string(payload),
		"projectId": updateProject.ID,
	}).Info("Updating project")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects/%s", c.BaseURL, *updateProject.ID), bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req = req.WithContext(ctx)
	res := Project{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Project": fmt.Sprintf("%+v", res),
	}).Debug("Project updated")
	return &res, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectId string) error {
	log.WithFields(log.Fields{
		"projectId": projectId,
	}).Info("Deleting project")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s", c.BaseURL, projectId), nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	if err := c.sendRequest(req, nil); err != nil {
		return err
	}
	return nil
}
