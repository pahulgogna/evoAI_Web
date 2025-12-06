package urlManager

import (
	"fmt"
	"strings"
)

func IsUrlValid(url string, root string) (bool, string) {
	url = strings.TrimSpace(url)

	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
		return true, url
	} 
	if strings.HasPrefix(url, "#") {
		return false, url
	}
	if strings.HasPrefix(url, "//") {
		return true, fmt.Sprintf("https:%s", url)
	}
	if strings.HasPrefix(url, "/") {
		if root == "" {
			return false, ""
		}
		return true, fmt.Sprintf("%s%s", root, url)
	}

	return false, url
}

func RankUrl(url string, query string, level int) int {

	rank := 0

	for _, word := range strings.Split(query, " ") {
		if strings.Contains(url, word) && !isFillerWord(word) {
			rank += 10
		}
	}

	rank += 10 / (level + 1) // reduces the rank, deeper the link is in the search.

	return rank
}

func isFillerWord(word string) bool {
	fillerWords := map[string]struct{}{
		"a": {}, "an": {}, "the": {}, "and": {}, "or": {}, "but": {}, "if": {}, "because": {},
		"as": {}, "what": {}, "where": {}, "when": {}, "how": {}, "who": {}, "which": {},
		"this": {}, "that": {}, "these": {}, "those": {}, "is": {}, "are": {}, "was": {},
		"were": {}, "be": {}, "been": {}, "being": {}, "have": {}, "has": {}, "had": {},
		"do": {}, "does": {}, "did": {}, "at": {}, "by": {}, "for": {}, "from": {}, "in": {},
		"into": {}, "of": {}, "off": {}, "on": {}, "onto": {}, "out": {}, "over": {}, "up": {},
		"with": {}, "to": {}, "am": {},
	}

	_, exists := fillerWords[strings.ToLower(word)]
	return exists
}