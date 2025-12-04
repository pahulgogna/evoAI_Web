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

type SearchRequest struct {
	Query            string `json:"query"`
	RequiredResults  int    `json:"required_results"`
	DnsAddress       string `json:"dns_address,omitempty"`
}