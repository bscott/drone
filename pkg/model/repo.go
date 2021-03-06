package model

import (
	"fmt"
	"time"
)

const (
	ScmGit = "git"
	ScmHg  = "hg"
	ScmSvn = "svn"
)

const (
	HostGithub    = "github.com"
	HostBitbucket = "bitbucket.org"
	HostGoogle    = "code.google.com"
	HostCustom    = "custom"
)

const (
	DefaultBranchGit = "master"
	DefaultBranchHg  = "default"
	DefaultBranchSvn = "trunk"
)

const (
	githubRepoPattern           = "git://github.com/%s/%s.git"
	githubRepoPatternPrivate    = "git@github.com:%s/%s.git"
	bitbucketRepoPattern        = "https://bitbucket.org/%s/%s.git"
	bitbucketRepoPatternPrivate = "git@bitbucket.org:%s/%s.git"
)

type Repo struct {
	ID int64 `meddler:"id,pk" json:"id"`

	// the full, canonical name of the repository, for example:
	// github.com/bradrydzewski/go.stripe
	Slug string `meddler:"slug" json:"slug"`

	// the hosting service where the repository is stored,
	// such as github.com, bitbucket.org, etc
	Host string `meddler:"host" json:"host"`

	// the owner of the repository on the host system.
	// for example, the Github username.
	Owner string `meddler:"owner" json:"owner"`

	// URL-friendly version of a repository name on the
	// host system.
	Name string `meddler:"name" json:"name"`

	// A value of True indicates the repository is closed source,
	// while a value of False indicates the project is open source.
	Private bool `meddler:"private" json:"private"`

	// A value of True indicates the repository is disabled and
	// no builds should be executed
	Disabled bool `meddler:"disabled" json:"disabled"`

	// A value of True indicates that pull requests are disabled
	// for the repository and no builds will be executed
	DisabledPullRequest bool `meddler:"disabled_pr" json:"disabled_pr"`

	// indicates the type of repository, such as
	// Git, Mercurial, Subversion or Bazaar.
	SCM string `meddler:"scm" json:"scm"`

	// the repository URL, for example:
	// git://github.com/bradrydzewski/go.stripe.git
	URL string `meddler:"url" json:"url"`

	// username and password requires to authenticate
	// to the repository
	Username string `meddler:"username" json:"username"`
	Password string `meddler:"password" json:"password"`

	// RSA key pair that will injected into the virtual machine
	// .ssh/id_rsa and .ssh/id_rsa.pub files.
	PublicKey  string `meddler:"public_key"  json:"public_key"`
	PrivateKey string `meddler:"private_key" json:"public_key"`

	// Parameters stored external to the repository in YAML
	// format, injected into the Build YAML at runtime.
	Params map[string]string `meddler:"params,gob" json:"-"`

	// the amount of time, in seconds the build will execute
	// before exceeding its timelimit and being killed.
	Timeout int64 `meddler:"timeout" json:"timeout"`

	// Indicates the build should be executed in priveleged
	// mode. This could, for example, be used to run Docker in Docker.
	Priveleged bool `meddler:"priveleged" json:"priveleged"`

	// Foreign keys signify the User that created
	// the repository and team account linked to
	// the repository.
	UserID int64 `meddler:"user_id"  json:"user_id"`
	TeamID int64 `meddler:"team_id"  json:"team_id"`

	Created time.Time `meddler:"created,utctime" json:"created"`
	Updated time.Time `meddler:"updated,utctime" json:"updated"`
}

// Creates a new repository
func NewRepo(host, owner, name, scm, url string) (*Repo, error) {
	repo := Repo{}
	repo.URL = url
	repo.SCM = scm
	repo.Host = host
	repo.Owner = owner
	repo.Name = name
	repo.Slug = fmt.Sprintf("%s/%s/%s", host, owner, name)
	key, err := generatePrivateKey()
	if err != nil {
		return nil, err
	}

	repo.PublicKey = marshalPublicKey(&key.PublicKey)
	repo.PrivateKey = marshalPrivateKey(key)
	return &repo, nil
}

// Creates a new GitHub repository
func NewGitHubRepo(owner, name string, private bool) (*Repo, error) {
	var url string
	switch private {
	case false:
		url = fmt.Sprintf(githubRepoPattern, owner, name)
	case true:
		url = fmt.Sprintf(githubRepoPatternPrivate, owner, name)
	}
	return NewRepo(HostGithub, owner, name, ScmGit, url)
}

// Creates a new Bitbucket repository
func NewBitbucketRepo(owner, name string, private bool) (*Repo, error) {
	var url string
	switch private {
	case false:
		url = fmt.Sprintf(bitbucketRepoPattern, owner, name)
	case true:
		url = fmt.Sprintf(bitbucketRepoPatternPrivate, owner, name)
	}
	return NewRepo(HostGithub, owner, name, ScmGit, url)
}

func (r *Repo) DefaultBranch() string {
	switch r.SCM {
	case ScmGit:
		return DefaultBranchGit
	case ScmHg:
		return DefaultBranchHg
	case ScmSvn:
		return DefaultBranchSvn
	default:
		return DefaultBranchGit
	}
}
