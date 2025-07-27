package htb

import (
	"fmt"
	"time"
)

// Challenge represents a HackTheBox challenge
type Challenge struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Difficulty  string   `json:"difficulty"`
	Points      int      `json:"points"`
	Solves      int      `json:"solves"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags,omitempty"`
	Released    string   `json:"released,omitempty"`
}

// Machine represents a HackTheBox machine
type Machine struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	OS         string  `json:"os"`
	Difficulty string  `json:"difficulty"`
	IPAddress  string  `json:"ip_address,omitempty"`
	Status     string  `json:"status"`
	UserOwned  bool    `json:"user_owned"`
	RootOwned  bool    `json:"root_owned"`
	Released   string  `json:"released,omitempty"`
	Rating     float64 `json:"rating,omitempty"`
	Active     bool    `json:"active"`
	Retired    bool    `json:"retired"`
	ExpiresAt  string  `json:"expires_at,omitempty"`
}

// User represents a HackTheBox user profile
type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Points         int    `json:"points"`
	Rank           string `json:"rank"`
	Subscription   string `json:"subscription"`
	SolvesCount    int    `json:"solves_count"`
	Country        string `json:"country,omitempty"`
	University     string `json:"university,omitempty"`
	CanAccessVIP   bool   `json:"canAccessVIP"`
	IsDedicatedVIP bool   `json:"isDedicatedVip"`
}

// SubmissionResult represents the result of a flag submission
type SubmissionResult struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	PointsAwarded int    `json:"points_awarded,omitempty"`
	FirstBlood    bool   `json:"first_blood,omitempty"`
}

// SearchResult represents search results from HTB API
type SearchResult struct {
	Machines   []SearchItem `json:"machines,omitempty"`
	Challenges []SearchItem `json:"challenges,omitempty"`
	Users      []SearchItem `json:"users,omitempty"`
}

// SearchItem represents a single search result item
type SearchItem struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// ChallengeListResponse represents the response from challenge list API
type ChallengeListResponse struct {
	Challenges []Challenge `json:"challenges"`
}

// MachineListResponse represents the response from machine list API
type MachineListResponse struct {
	Data []Machine `json:"data"`
}

// UserInfoResponse represents the response from user info API
type UserInfoResponse struct {
	Info User `json:"info"`
}

// ActiveMachineResponse represents the response from active machine API
type ActiveMachineResponse struct {
	Info *Machine `json:"info"`
}

// FlagSubmissionRequest represents a flag submission request
type FlagSubmissionRequest struct {
	Flag        string `json:"flag"`
	ID          int    `json:"id,omitempty"`
	ChallengeID string `json:"challenge_id,omitempty"`
	Difficulty  string `json:"difficulty,omitempty"`
}

// MachineActionRequest represents a machine action request (start/stop)
type MachineActionRequest struct {
	MachineID int `json:"machine_id,omitempty"`
	ID        int `json:"id,omitempty"`
}

// PaginatedRequest represents common pagination parameters
type PaginatedRequest struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
}

// ChallengeFilter represents filters for challenge listing
type ChallengeFilter struct {
	Category   string `json:"category,omitempty"`
	Difficulty string `json:"difficulty,omitempty"`
	Status     string `json:"status,omitempty"` // active, retired
	SortBy     string `json:"sort_by,omitempty"`
	PaginatedRequest
}

// MachineFilter represents filters for machine listing
type MachineFilter struct {
	OS         string `json:"os,omitempty"`
	Difficulty string `json:"difficulty,omitempty"`
	Status     string `json:"status,omitempty"` // active, retired
	SortBy     string `json:"sort_by,omitempty"`
	PaginatedRequest
}

// ServerStatus represents the MCP server health status
type ServerStatus struct {
	Status       string    `json:"status"`
	Version      string    `json:"version"`
	HTBAPIStatus string    `json:"htb_api_status"`
	Uptime       string    `json:"uptime"`
	Timestamp    time.Time `json:"timestamp"`
}

// Error represents an API error response
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// HTBAPIError represents an error from the HTB API
type HTBAPIError struct {
	StatusCode int
	Message    string
	Response   string
}

func (e *HTBAPIError) Error() string {
	return fmt.Sprintf("HTB API error (status %d): %s", e.StatusCode, e.Message)
}

// DifficultyLevel represents the difficulty levels used by HTB
type DifficultyLevel string

const (
	DifficultyEasy   DifficultyLevel = "Easy"
	DifficultyMedium DifficultyLevel = "Medium"
	DifficultyHard   DifficultyLevel = "Hard"
	DifficultyInsane DifficultyLevel = "Insane"
)

// MachineType represents the type of machine
type MachineType string

const (
	MachineTypeActive  MachineType = "active"
	MachineTypeRetired MachineType = "retired"
	MachineTypeRelease MachineType = "release"
)

// SubscriptionType represents user subscription levels
type SubscriptionType string

const (
	SubscriptionFree    SubscriptionType = "free"
	SubscriptionVIP     SubscriptionType = "vip"
	SubscriptionVIPPlus SubscriptionType = "vip+"
)
