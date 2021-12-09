package main

import (
	"time"
)

type RunTaskRequest struct {
	PayloadVersion             int       `json:"payload_version"`
	AccessToken                string    `json:"access_token"`
	TaskResultID               string    `json:"task_result_id"`
	TaskResultEnforcementLevel string    `json:"task_result_enforcement_level"`
	TaskResultCallbackUrl      string    `json:"task_result_callback_url"`
	RunAppUrl                  string    `json:"run_app_url"`
	RunID                      string    `json:"run_id"`
	RunMessage                 string    `json:"run_message"`
	RunCreatedAt               time.Time `json:"run_created_at"`
	RunCreatedBy               string    `json:"run_created_by"`
	WorkspaceID                string    `json:"workspace_id"`
	WorkspaceName              string    `json:"workspace_name"`
	WorkspaceAppUrl            string    `json:"workspace_app_url"`
	OrganizationName           string    `json:"organization_name"`
	PlanJsonApiUrl             string    `json:"plan_json_api_url"`
	VcsRepoUrl                 string    `json:"vcs_repo_url"`
	VcsBranch                  string    `json:"vcs_branch"`
	VcsPullRequestUrl          string    `json:"vcs_pull_request_url,omitempty"`
	VcsCommitUrl               string    `json:"vcs_commit_url"`
}

type RunTaskResponse struct {
	ID      string `jsonapi:"primary,task-results"`
	Status  string `jsonapi:"attr,status"`
	Message string `jsonapi:"attr,message"`
	Url     string `jsonapi:"attr,url"`
}
