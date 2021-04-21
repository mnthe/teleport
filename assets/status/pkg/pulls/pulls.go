package pulls

import (
	"context"
	"time"

	"github.com/gravitational/teleport/assets/statusctl/pkg/constants"

	"github.com/google/go-github/v35/github"
	"github.com/gravitational/trace"
	"golang.org/x/oauth2"
)

type Config struct {
	AccessToken string
	Verbose     bool
}

func (c *Config) CheckAndSetDefaults() error {
	if c.Verbose && c.AccessToken == "" {
		return trace.BadParameter("verbose unavailable without an access token")
	}

	return nil
}

func Print(ctx context.Context, c *Config) error {
	prs, err := fetch(ctx, c)
	if err != nil {
		return trace.Wrap(err)
	}

	if err := display(ctx, prs); err != nil {
		return trace.Wrap(err)
	}

	return nil
}

type pullRequest struct {
	//PR      *github.PullRequest
	//Reviews []*github.PullRequestReview

	// group is how we categorize PRs. A few examples, "code", "rfd", "docs",
	// "draft", "backport".
	group string

	// openFor is how long the PR has been open.
	openFor time.Duration

	// number is the GitHub PR number, like #1234.
	number int

	// author is the GitHub handle of the PR author.
	author string

	// title is the title of the PR.
	title string

	// approvers is a slice of GitHub handles that have approved the PR. Only
	// available in verbose mode.
	approvers []string
}

func fetch(ctx context.Context, c *Config) ([]pullRequest, error) {
	if err := c.CheckAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}
	client := newClient(ctx, c.AccessToken)

	var prs []pullRequest

	// ?
	popts := &github.PullRequestListOptions{
		State: constants.Open,
		ListOptions: github.ListOptions{
			PerPage: constants.PageSize,
		},
	}

	// Paginate and get all PRs.
	for {
		page, resp, err := client.PullRequests.List(ctx,
			constants.Organization, constants.Repository, popts)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		for _, pr := range page {
			ropts := &github.ListOptions{
				PerPage: 20,
			}
			reviews, _, err := client.PullRequests.ListReviews(context.Background(), "gravitational", "teleport", pr.GetNumber(), ropts)
			if err != nil {
				return nil, err
			}

			pulls = append(pulls, &pullRequest{
				pr:      pr,
				reviews: reviews,
			})
		}

		if resp.NextPage == 0 {
			break
		}
		ropts.Page = resp.NextPage
	}

	return pulls, nil
}

func display(ctx context.Context, pr []pullRequest) error {
	return nil
}

//func fetchm() ([]*pullRequest, error) {
//	var pulls []*pullRequest
//
//	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
//		&oauth2.Token{AccessToken: ""},
//	))
//	client := github.NewClient(tc)
//
//	iopts := &github.IssueListByRepoOptions{
//		Milestone: "55",
//		State:     "open",
//	}
//	milestone, resp, err := client.Issues.ListByRepo(context.Background(), "gravitational", "teleport", iopts)
//	if err != nil {
//		return nil, err
//	}
//
//	for {
//		page, resp, err := client.Issues.GetMilestone(context.Background(), "gravitational", "teleport", 55)
//		if err != nil {
//			return nil, err
//		}
//
//		for _, pr := range page {
//			lopts := &github.ListOptions{
//				PerPage: 20,
//			}
//			reviews, _, err := client.PullRequests.ListReviews(context.Background(), "gravitational", "teleport", pr.GetNumber(), lopts)
//			if err != nil {
//				return nil, err
//			}
//
//			pulls = append(pulls, &pullRequest{
//				pr:      pr,
//				reviews: reviews,
//			})
//		}
//
//		if resp.NextPage == 0 {
//			break
//		}
//		ropts.Page = resp.NextPage
//	}
//
//	return pulls, nil
//}

//func exclude(pr *github.PullRequest) bool {
//	if pr.GetState() == "draft" {
//		return true
//	}
//	for _, label := range pr.Labels {
//		if *label.Name == "documentation" {
//			return true
//		}
//	}
//	return false
//}
//func isMember(team string, name string) bool {
//	for _, s := range teams[team] {
//		if name == s {
//			return true
//		}
//	}
//	return false
//}
//
//func isAnyTeam(name string) bool {
//	for _, v := range teams {
//		for _, vv := range v {
//			if vv == name {
//				return true
//			}
//		}
//	}
//	return false
//}
//
//func printTeam(team string, prs []*github.PullRequest) {
//	n := 0
//	groups := map[string][]*github.PullRequest{}
//
//	for _, pr := range prs {
//		user := *pr.GetUser().Login
//
//		if exclude(pr) {
//			continue
//		}
//		if !isMember(team, user) {
//			continue
//		}
//
//		n = n + 1
//
//		var ok bool
//		pullslice := []*github.PullRequest{}
//
//		if pullslice, ok = groups[user]; ok {
//			pullslice = groups[user]
//		}
//		pullslice = append(pullslice, pr)
//		groups[user] = pullslice
//	}
//
//	if n == 0 {
//		return
//	}
//
//	fmt.Printf("--------------------------------------------------------------------------------\n")
//	fmt.Printf("Team: %v, Open: %v\n", team, n)
//	fmt.Printf("--------------------------------------------------------------------------------\n")
//
//	for k, v := range groups {
//		for _, vv := range v {
//			duration := time.Now().Sub(*vv.CreatedAt)
//			humanDuration := fmt.Sprintf("%v", duration.Round(24*time.Hour))
//
//			fmt.Printf("%-5v %-20v %-10v %v.\n", *vv.Number, k, humanDuration, *vv.Title)
//		}
//	}
//}

//type summaryView struct {
//	// draft, docs, rfd, code
//	category string
//	openfor  string
//	team     string
//	number   int
//	author   string
//	title    string
//	count    int
//}
//
//func summary(pulls []*pullRequest) {
//	sv := make([]summaryView, 0, len(pulls))
//
//	var cn int
//	var rn int
//	var dn int
//
//	for _, pull := range pulls {
//		var n int
//		for _, review := range pull.reviews {
//			if review.GetState() == "APPROVED" {
//				n += 1
//			}
//		}
//
//		switch getCategory(pull.pr) {
//		case "code":
//			cn += 1
//		case "docs":
//			dn += 1
//		case "rfd":
//			rn += 1
//		}
//
//		team := "external"
//		if n, ok := teams[pull.pr.GetUser().GetLogin()]; ok {
//			team = n
//		}
//
//		duration := time.Now().Sub(pull.pr.GetCreatedAt())
//		humanDuration := fmt.Sprintf("%vd", math.Ceil(duration.Hours()/24))
//
//		sv = append(sv, summaryView{
//			category: getCategory(pull.pr),
//			team:     team,
//			openfor:  humanDuration,
//			number:   pull.pr.GetNumber(),
//			author:   pull.pr.GetUser().GetLogin(),
//			count:    n,
//			title:    pull.pr.GetTitle(),
//		})
//	}
//
//	sort.Slice(sv, func(i, j int) bool {
//		if sv[i].team < sv[j].team {
//			return true
//		}
//		if sv[i].team > sv[j].team {
//			return false
//		}
//		return sv[i].author < sv[j].author
//	})
//
//	fmt.Printf("code: open %v\n", cn)
//	printSummary(sv, "code")
//
//	fmt.Printf("\nrfd: open: %v\n", rn)
//	printSummary(sv, "rfd")
//
//	fmt.Printf("\ndocs: open: %v\n", dn)
//	printSummary(sv, "docs")
//}
//
//func printSummary(sv []summaryView, category string) {
//	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
//	for _, v := range sv {
//		if v.category != category {
//			continue
//		}
//		fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t", v.number, v.count, v.openfor, v.author, v.team, v.title))
//	}
//	w.Flush()
//
//}
//
//func getCategory(pr *github.PullRequest) string {
//	if pr.GetDraft() {
//		return "draft"
//	}
//	if hasLabel(pr, "documentation") {
//		return "docs"
//	}
//	if hasLabel(pr, "rfd") {
//		return "rfd"
//	}
//	return "code"
//}
//
//func isDraft(pr *github.PullRequest) bool {
//	fmt.Printf("%v.\n", pr.GetState())
//	if pr.GetState() == "draft" {
//		return true
//	}
//	return false
//}
//
//func hasLabel(pr *github.PullRequest, name string) bool {
//	for _, label := range pr.Labels {
//		if label.GetName() == name {
//			return true
//		}
//	}
//	return false
//}

//func main() {
//	token, err := config.ReadToken()
//	if err != nil {
//		log.Fatalf("Failed to read in GitHub OAuth2 token: %v.", err)
//	}
//
//	//pulls, err := fetch()
//	//if err != nil {
//	//	log.Fatalf("Failed to fetch: %v.", err)
//	//}
//
//	////if len(os.Args) > 1 && os.Args[1] == "summary" {
//	////	summary(pulls)
//	////} else {
//	////	milestone()
//	////}
//
//	//printTeam("security", prs)
//	//printTeam("scale", prs)
//	//printTeam("sshkube", prs)
//	//printTeam("appdb", prs)
//	//printTeam("release", prs)
//
//	//printTeam("", prs)
//}

func newClient(ctx context.Context, accessToken string) *github.Client {
	// If no accessToken was passed in, use a (rate limited) unauthenticated.
	if accessToken == "" {
		return github.NewClient(nil)
	}

	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	))
	return github.NewClient(tc)
}