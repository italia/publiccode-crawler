package crawler

var (
	clientApis map[string]func(domain Domain, url string, repositories chan Repository) (string, error)
)

func RegisterClientApis() {
	clientApis = make(map[string]func(domain Domain, url string, repositories chan Repository) (string, error))

	clientApis["bitbucket"] = RegisterBitbucketAPI()
	clientApis["github"] = RegisterGithubAPI()
	clientApis["gitlab"] = RegisterGitlabAPI()
}

func GetClientApiCrawler(clientApi string) func(domain Domain, url string, repositories chan Repository) (string, error) {
	return clientApis[clientApi]
}
