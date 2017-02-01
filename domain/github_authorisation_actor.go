package domain

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubAuthorisationActor struct {
	*AWSActor
	c                chan string
	OAuthConfig      *oauth2.Config
	oauthStateString string
}

func (g *GithubAuthorisationActor) Run(interface{}) error {
	go g.webServer()

	token, err := g.retrieveToken()
	if err != nil {
		return err
	}
	g.Context.Props["ApplicationRepositoryOAuthToken"] = token
	return nil
}

func (g *GithubAuthorisationActor) retrieveToken() (string, error) {
	g.OAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/GithubCallback",
		ClientID:     "b869de5987f281f73339",
		ClientSecret: "f3469cd505ec31d72dd8f003285eb609558fd79f",
		Scopes:       []string{"repo"},
		Endpoint:     github.Endpoint,
	}
	// Some random string, random for each request
	g.oauthStateString = "random"
	g.c = make(chan string, 1)
	err := exec.Command("open", "http://localhost:3000/GithubLogin").Start()
	if err != nil {
		return "", err
	}

	token := <-g.c
	close(g.c)
	return token, nil

}

func (g *GithubAuthorisationActor) webServer() {
	http.HandleFunc("/GithubLogin", g.handleGithubLogin)
	http.HandleFunc("/GithubCallback", g.handleGithubCallback)
	fmt.Println(http.ListenAndServe(":3000", nil))
}

func (g *GithubAuthorisationActor) handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	url := g.OAuthConfig.AuthCodeURL(g.oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (g *GithubAuthorisationActor) handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != g.oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", g.oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := g.OAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	g.c <- token.AccessToken
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Authorised with Github, you can now close this tab!")
}
