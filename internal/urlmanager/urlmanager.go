package urlmanager

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type UrlManager struct {
	baseURL     *url.URL
	baseDomain  string
	visitedURLs map[string]bool
	mutex       sync.RWMutex
}

// New creates a new UrlManager
func New(startURL string) (*UrlManager, error) {
	// Parse the start URL
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		return nil, err
	}

	// ensure the URL has a scheme (default to https)
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	// Extract domin (hostname without port)
	baseDomain := parsedURL.Hostname()
	if baseDomain == "" {
		return nil, fmt.Errorf("invalid URL: %s", startURL)
	}

	// Initialize the UrlManager
	return &UrlManager{
		baseURL:     parsedURL,
		baseDomain:  baseDomain,
		visitedURLs: make(map[string]bool),
	}, nil
}

// Checks if the url is within the base domain
func (u *UrlManager) isWithinBaseDomain(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	hostname := parsedURL.Hostname()
	if hostname == "" {
		return false
	}

	return hostname == u.baseDomain
}

// Normalize the url - parse, remove fragment, and normalize path
// baseURLStr is the URL of the page where this link was found (for resolving relative URLs)
func (u *UrlManager) normalizeUrl(urlStr string, baseURLStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	if parsedURL.Scheme == "" {
		// Use the base URL of the current page if provided, otherwise use the initial base URL
		var baseURL *url.URL
		if baseURLStr != "" {
			baseURL, err = url.Parse(baseURLStr)
			if err != nil {
				return "", err
			}
		} else {
			baseURL = u.baseURL
		}
		parsedURL = baseURL.ResolveReference(parsedURL)
	}

	parsedURL.Fragment = ""
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
	return parsedURL.String(), nil
}

// ShouldVisit checks if a URL should be visited (not already visited and within base domain)
// baseURLStr is the URL of the page where this link was found (for resolving relative URLs)
func (u *UrlManager) ShouldVisit(urlStr string, baseURLStr string) (bool, string, error) {
	normalizedUrl, err := u.normalizeUrl(urlStr, baseURLStr)
	if err != nil {
		return false, "", err
	}

	sameDomain := u.isWithinBaseDomain(normalizedUrl)
	if !sameDomain {
		return false, "", nil
	}

	u.mutex.RLock()
	visited := u.visitedURLs[normalizedUrl]
	u.mutex.RUnlock()

	return !visited, normalizedUrl, nil
}

func (u *UrlManager) MarkAsVisited(urlStr string) error {
	// When marking as visited, we already have the normalized URL, so no need for base URL
	normalizedUrl, err := u.normalizeUrl(urlStr, "")
	if err != nil {
		return err
	}

	u.mutex.Lock()
	u.visitedURLs[normalizedUrl] = true
	u.mutex.Unlock()
	return nil
}
