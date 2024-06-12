package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

var (
	client *github.Client
)

func init() {
	// Inicializa o cliente GitHub
	client = newGitHubClient()
}

func newGitHubClient() *github.Client {
	// Obtém o token de acesso pessoal do ambiente
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("Token de acesso pessoal não configurado")
	}

	ctx := context.Background()

	// Cria um cliente de autenticação com o token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

type CreateRepoRequest struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

func CreateRepository(w http.ResponseWriter, r *http.Request) {
	var req CreateRepoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if req.Name == "" {
		http.Error(w, "Repository name is required", http.StatusBadRequest)
		return
	}

	owner := "RubenCsantos"
	repo := &github.Repository{
		Name:    github.String(req.Name),
		Private: github.Bool(req.Private),
	}

	// Cria o repositório
	_, _, err := client.Repositories.Create(context.Background(), owner, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Repositório criado com sucesso!")
}

func DeleteRepository(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner := vars["owner"]
	repoName := vars["repo"]

	if owner == "" || repoName == "" {
		http.Error(w, "Owner and repository name are required", http.StatusBadRequest)
		return
	}

	// Apaga o repositório
	_, err := client.Repositories.Delete(context.Background(), owner, repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Repositório excluído com sucesso!")
}

func ListRepositories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	if user == "" {
		http.Error(w, "User is required", http.StatusBadRequest)
		return
	}

	// Lista os repositórios do usuário
	repos, _, err := client.Repositories.List(context.Background(), user, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	for _, repo := range repos {
		fmt.Fprintf(w, "%s\n", *repo.FullName)
	}
}

func ListPullRequests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner := vars["owner"]
	repoName := vars["repo"]

	if owner == "" || repoName == "" {
		http.Error(w, "Owner and repository name are required", http.StatusBadRequest)
		return
	}

	// Lista os pull requests do repo
	prs, _, err := client.PullRequests.List(context.Background(), owner, repoName, &github.PullRequestListOptions{State: "open"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	for _, pr := range prs {
		fmt.Fprintf(w, "%d: %s\n", *pr.Number, *pr.Title)
	}
}
