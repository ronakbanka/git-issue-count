package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/stackimpact/stackimpact-go"

)

type issue struct {
	Data template.HTML
}

func main() {

	log.Printf("Server started listening on %s", os.Getenv("PORT"))
        
        // Add stackimpact analysis
        stackkey := os.Getenv("STACKIMPACT_KEY")
        agent := stackimpact.NewAgent()
        agent.Configure(stackkey, "github-issue-count")

	http.HandleFunc("/", index)
	http.HandleFunc("/getinfo", getinfo)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// get repo name
		fmt.Println("github repo URL:", r.Form["repo"])

	}
}

// handler func to return formatted output
func getinfo(w http.ResponseWriter, r *http.Request) {

	var data string

	git_repo := r.FormValue("repo")

	u, err := url.Parse(git_repo)
	if err != nil {
		log.Fatal(err)
	}
	if u.Host == "github.com" {
		details := strings.Split(u.Path, "/")
		owner := details[1]
		repo := details[2]

		log.Printf("checking info for %s", git_repo)

		client := github.NewClient(nil)

		open_issues, err := hasOpenIssues(client, owner, repo) // Check if there is any open issues
		if err != nil {
			log.Fatal(err)
		} else if open_issues == 0 {
			data = fmt.Sprint("<b>There are no open issues</b>")
		} else {
			num_1d, num_7d, num_rest, total_issues, _ := getIssueCount(client, owner, repo)
			data = fmt.Sprintf("Github Repo: https://github.com/%s/%s <br>"+
				"Total number of open issues : %d <br>"+
				"Number of open issues in last 24hr : %d <br>"+
				"Number of open issues that were opened more than 24 hours ago but less than 7 days ago : %d <br>"+
				"Number of open issues that were opened more than 7 days ago : %d <br>",
				owner, repo, total_issues, num_1d, num_7d, num_rest)
		}
	} else {
		data = fmt.Sprint("<b>Enter a valid Github url</b>")
	}

	t, _ := template.ParseFiles("index.gtpl")
	info := issue{Data: template.HTML(data)}
	t.Execute(w, &info)
}
