package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	. "github.com/drone/drone/pkg/model"
	"github.com/russross/meddler"
)

const commitRecentTemplate = `
SELECT r.slug, r.host, r.owner, r.name,
c.status, c.started, c.finished, c.duration, c.hash, c.branch, c.pull_request,
c.author, c.gravatar, c.timestamp, c.message, c.created, c.updated
FROM repos r, commits c
WHERE r.id = c.repo_id
AND r.slug IN (%s)
ORDER BY c.created desc
LIMIT 20
`

const teamCommitRecentStmt = `
SELECT r.slug, r.host, r.owner, r.name,
c.status, c.started, c.finished, c.duration, c.hash, c.branch, c.pull_request,
c.author, c.gravatar, c.timestamp, c.message, c.created, c.updated
FROM repos r, commits c
WHERE r.team_id = ?
AND   r.id = c.repo_id
AND   c.status IN ('Success', 'Failure')
ORDER BY c.created desc
LIMIT 10
`

var (
	db          *sql.DB
	repoList    []interface{}
	commitsStmt string
)

func reposQuery() string {
	list := strings.Split(*repos, ",")

	repoList = make([]interface{}, len(list))
	for i, v := range list {
		repoList[i] = interface{}(v)
	}

	var params bytes.Buffer
	for i := range repoList {
		if i < len(repoList)-1 {
			params.WriteString(fmt.Sprintf("?, "))
		} else {
			params.WriteString(fmt.Sprintf("?"))
		}
	}

	return fmt.Sprintf(commitRecentTemplate, params.String())
}

func listWallCommits() ([]*RepoCommit, error) {
	var commits []*RepoCommit
    var err error

	if *teams != "" {
        err = meddler.QueryAll(db, &commits, teamCommitRecentStmt, *team)
	} else if *repos != "" {
        commits := setupReposQuery()
        err = meddler.QueryAll(db, &commits, commits, repoList...)
	}

	return commits, err
}
