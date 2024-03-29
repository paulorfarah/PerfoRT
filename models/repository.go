package models

import (
	//	"fmt"
	//	"context"
	"gorm.io/gorm"
	//	"github.com/shurcool/githubv4"
	//"golang.org/x/oauth2"
)

type Repository struct {
	Model
	PlatformID   uint // `gorm:"unique_index:idx_repository"`
	Platform     Platform
	Name         string `gorm:"not null;unique_index:idx_repository;size:2048"`
	Description  string `gorm:"size:2048"`
	IsPrivate    bool   `gorm:"not null" sql:"DEFAULT:false"`
	Commits      []Commit
	Measurements []Measurement
}

func (r *Repository) TableName() string {
	return "repositories"
}

func CreateRepository(db *gorm.DB, repository *Repository) (uint, error) {
	err := db.Create(repository).Error
	if err != nil {
		return 0, err
	}
	return repository.ID, nil
}

func FindRepositoryByName(db *gorm.DB, name string) (*Repository, error) {
	var repository Repository
	res := db.Where("name = ?", name).First(&repository)
	return &repository, res.Error
}

//func (r *Repository) UpdateCommits(db *gorm.DB, branchHash xxx, since time, until time)  {
//	fmt.Println("TODO")
//}

/*
func (r *Repository)  Issues() error {
	var query struct {
		Repository struct {
			Nodes []struct{
				Issue struct {
					Number	githubv4.Int
					Repository struct {
						NameWithOwner githubv4.String
					}
				}`graphql:"... on Issue"`
			}
		}`graphql:"searc\h(first: 100, query: $searchQuery, type: $searchType)"`

	}

	variables := map[string]interface{}{
		"searchQuery": githubv4.String("repo:cockroachdb/cockroach state: open teamcity: failed tests on master"),
		"searchType":  githubv4.SearchTypeIssue,
	}

	client := GetClient()
	// query issues
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println("error querying issues")
	}
//	fmt.Println(query.Search.Nodes.Issue.)


	return nil
}
*/
