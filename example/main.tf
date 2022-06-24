terraform {
  required_providers {
    sodium = {
      source  = "github.com/sodiumprovider/sodium"
    }
    github = {
      source  = "integrations/github"
      version = ">= 4.5.2"
    }
  }
}

provider "github" {
  owner = "org_name"
  token = "github_token"
}


# To make sure the repository exists and the correct permissions are set.
data "github_repository" "main" {
  full_name = "org_name/repo_name"
}

data "github_actions_public_key" "gh_actions_public_key" {
  repository = "repo_name"
}

data "sodium_encrypted_item" "foo" {
    public_key = data.github_actions_public_key.gh_actions_public_key.key
    content_base64 = base64encode("SuperSecretPassword")
}

resource "github_actions_secret" "gh_actions_secret" {
  repository       = "repo_name"
  secret_name      = "SECRET_FOO"
  encrypted_value  = data.sodium_encrypted_item.foo.encrypted_value_base64
}