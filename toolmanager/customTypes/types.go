package customtypes

type Snippet struct {
	Description string
	Language string
	Code string
	Dependencies []string
}

type Store struct {
	Tools map[string]Snippet `yaml:"tools"`
}

type CreateRequestSchema struct {
	Name string
	Tool Snippet
}
