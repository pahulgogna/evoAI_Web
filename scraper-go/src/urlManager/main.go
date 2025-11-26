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



