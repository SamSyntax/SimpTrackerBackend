package utils

import "sync"

type Keyword struct {
	ID      int32  `json:"id"`
	Keyword string `json:"keyword"`
	Active  bool   `json:"active"`
	Count   int64  `json:"total_count"`
}

var (
	activeKeywords     []string
	activeKeywordsLock sync.RWMutex
)

// SetActiveKeywords initializes the active keywords list
func SetActiveKeywords(keywords []string) {
	activeKeywordsLock.Lock()
	defer activeKeywordsLock.Unlock()
	activeKeywords = keywords
}

// GetActiveKeywords retrieves a copy of the active keywords list
func GetActiveKeywords() []string {
	activeKeywordsLock.RLock()
	defer activeKeywordsLock.RUnlock()
	keywordsCopy := make([]string, len(activeKeywords))
	copy(keywordsCopy, activeKeywords)
	return keywordsCopy
}

// AddKeyword adds a keyword to the active keywords list
func AddKeyword(keyword string) {
	activeKeywordsLock.Lock()
	defer activeKeywordsLock.Unlock()
	activeKeywords = append(activeKeywords, keyword)
}

// RemoveKeyword removes a keyword from the active keywords list
func RemoveKeyword(keyword string) {
	activeKeywordsLock.Lock()
	defer activeKeywordsLock.Unlock()
	for i, kw := range activeKeywords {
		if kw == keyword {
			activeKeywords = append(activeKeywords[:i], activeKeywords[i+1:]...)
			break
		}
	}
}

// UpdateKeywords replaces the entire active keywords list
func UpdateKeywords(keywords []string) {
	activeKeywordsLock.Lock()
	defer activeKeywordsLock.Unlock()
	activeKeywords = keywords
}
