package models

import(
	"context"
	
	"golang.org/x/oauth2"
	"github.com/shurcool/githubv4"
)

//var loginStr struct {
//	Viewer struct{
//		Login	string
//		CreatedAt	time.Time
//		IsBountyHunter	bool
//		WebsiteURL	string
//	}
//}

func GetClient() *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "0cd723e75e96e3e4f20994ba9d996494f448b107"},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	return client
	// query login
//	err := client.Query(context.Background(), &loginStr, nil)
//	if err != nil {
//		fmt.Println("login error: ", err)
//	}
//	fmt.Println("Login: ", loginStr.Viewer.Login)
//	fmt.Println("CreatedAt: ", login.Viewer.CreatedAt)

}
