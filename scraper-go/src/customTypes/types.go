package customTypes

type Page struct {
	Source 	string // url
	Body 	string  // page body
}

type StoreUrl struct {
	Priority int
	Url 	 string
	Index 	 int
	Level 	 int
}
