package models

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcool/githubv4"
	"gorm.io/gorm"
)

type Issue struct {
	Model
	RepositoryID     int // uint `gorm:"not null;unique_index:idx_issue"`
	Repository       Repository
	Number           int  `gorm:"not null"` //issue number
	CreatedViaEmail  bool `gorm:"default:false"`
	PublishedAt      time.Time
	Title            string `gorm:"not null"`
	Author           uint   `gorm:"not null"`
	BodyText         string `gorm:"not null"`
	State            string `gorm:"not null"`
	Closed           bool   `gorm:"default:false"`
	ClosedAt         time.Time
	Editor           uint
	LastEditedAt     time.Time
	Locked           bool `gorm:"default:false"`
	ActiveLockReason string
	ResourcePath     string `gorm:"not null"`
	// UpdatedAt        time.Time
	Url string

	//milestone
	//labels
	//assignee
	//comments
	//participants
	//projectCards
	//Reactions
	//TimelineItems
	//UserContentEdits
}

func (i *Issue) TableName() string {
	return "issues"
}

func CreateIssue(db *gorm.DB, issue *Issue) (uint, error) {
	err := db.Create(issue).Error
	if err != nil {
		return 0, err
	}
	return issue.ID, nil
}

func FindIssueByRepository(db *gorm.DB, repository uint) (*Issue, error) {
	var issue Issue
	last := db.Last(&issue)
	if last.Error != nil {
		return nil, last.Error
	}
	return &issue, nil
}

//type IssueGH struct {
//	Number uint
//}

func GetIssues(lastIssue *Issue) []Issue {
	owner := "paulorfarah"
	name := "refactoring-python-code"
	fmt.Println("Read Issues")

	if lastIssue != nil {
		//since := githubv4.DateTime(lastIssue.IssueCreatedAt)
		var q struct {
			Repository struct {
				Issues struct {
					Nodes    []Issue
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"issues(first: 100, after:$issuesCursor, filterBy:{since: $since})"`
			} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
		}
		variables := map[string]interface{}{
			"repositoryOwner": githubv4.String(owner),
			"repositoryName":  githubv4.String(name),
			"since":           lastIssue.CreatedAt.String(),
			"issuesCursor":    (*githubv4.String)(nil),
		}
		client := GetClient()
		var allIssues []Issue
		for {
			err := client.Query(context.Background(), &q, variables)
			if err != nil {
				fmt.Println(err)
			}
			allIssues = append(allIssues, q.Repository.Issues.Nodes...)
			if !q.Repository.Issues.PageInfo.HasNextPage {
				break
			}
			variables["issuesCursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)
			fmt.Println(allIssues)
		}
		fmt.Println("Read", allIssues)
		//return allIssuesi
		return nil
	} else {
		fmt.Println("else")
		var q struct {
			Repository struct {
				Issues struct {
					Nodes    []Issue
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"issues(first: 100, after:$issuesCursor)"`
			} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
		}
		fmt.Println("q")
		variables := map[string]interface{}{
			"repositoryOwner": githubv4.String(owner),
			"repositoryName":  githubv4.String(name),
			"issuesCursor":    (*githubv4.String)(nil),
		}
		fmt.Println("variables")
		client := GetClient()
		var allIssues []Issue
		for {
			err := client.Query(context.Background(), &q, variables)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%+v\n", q.Repository.Issues.Nodes)
			allIssues = append(allIssues, q.Repository.Issues.Nodes...)
			fmt.Println("allIssues")
			if !q.Repository.Issues.PageInfo.HasNextPage {
				break
			}
			fmt.Println("break")
			variables["issuesCursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)
			fmt.Println("EndCursor")
			fmt.Println(allIssues)
		}
		return allIssues
	}
}
