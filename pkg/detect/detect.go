package detect

import "os"

type ProjectType string

const (
	Node    ProjectType = "Node"
	Go      ProjectType = "Go"
	Python  ProjectType = "Python"
	Rust    ProjectType = "Rust"
	PHP     ProjectType = "PHP"
	Ruby    ProjectType = "Ruby"
	Unknown ProjectType = "Unknown"
)

func Detect() ProjectType {
	switch {
	case exists("package.json"):
		return Node
	case exists("go.mod"):
		return Go
	case exists("requirements.txt"), exists("setup.py"), exists("pyproject.toml"):
		return Python
	case exists("Cargo.toml"):
		return Rust
	case exists("composer.json"):
		return PHP
	case exists("Gemfile"):
		return Ruby
	default:
		return Unknown
	}
}

func exists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

func GetExampleVars(pt ProjectType) map[string]string {
	base := map[string]string{
		"APP_ENV": "development",
	}
	switch pt {
	case Node:
		base["PORT"] = "3000"
		base["DATABASE_URL"] = "postgresql://localhost:5432/mydb"
		base["JWT_SECRET"] = "change_me"
		base["REDIS_URL"] = "redis://localhost:6379"
	case Go:
		base["PORT"] = "8080"
		base["DB_HOST"] = "localhost"
		base["DB_USER"] = "postgres"
		base["DB_NAME"] = "mydb"
		base["DB_PASSWORD"] = "change_me"
	case Python:
		base["PORT"] = "5000"
		base["DATABASE_URL"] = "postgresql://localhost:5432/mydb"
		base["SECRET_KEY"] = "change_me"
		base["DEBUG"] = "true"
	case Rust:
		base["PORT"] = "8080"
		base["DATABASE_URL"] = "postgresql://localhost:5432/mydb"
		base["RUST_LOG"] = "info"
	case PHP:
		base["DB_HOST"] = "localhost"
		base["DB_DATABASE"] = "mydb"
		base["DB_USERNAME"] = "root"
		base["DB_PASSWORD"] = "change_me"
		base["APP_KEY"] = "change_me"
	case Ruby:
		base["DATABASE_URL"] = "postgresql://localhost:5432/mydb"
		base["SECRET_KEY_BASE"] = "change_me"
		base["RAILS_ENV"] = "development"
	default:
		base["API_KEY"] = "your_api_key_here"
		base["DEBUG"] = "true"
	}
	return base
}
