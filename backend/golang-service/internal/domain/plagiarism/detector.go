package plagiarism

import (
	"context"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// Detector defines the plagiarism detection interface
type Detector interface {
	Compare(ctx context.Context, text1, text2 string) (float64, []Match, error)
	CheckAgainstCorpus(ctx context.Context, text string, submissionIDs []uuid.UUID) (*ReportDetails, error)
}

type detector struct {
	nGramSize int
}

// NewDetector creates a new plagiarism detector
func NewDetector(nGramSize int) Detector {
	if nGramSize < 2 {
		nGramSize = 5 // default to 5-grams
	}
	return &detector{nGramSize: nGramSize}
}

// Compare compares two texts and returns similarity score and matches
func (d *detector) Compare(ctx context.Context, text1, text2 string) (float64, []Match, error) {
	// Normalize texts
	norm1 := normalize(text1)
	norm2 := normalize(text2)

	// Generate n-grams
	ngrams1 := generateNGrams(norm1, d.nGramSize)
	ngrams2 := generateNGrams(norm2, d.nGramSize)

	// Build set for text2 ngrams
	set2 := make(map[string]bool)
	for _, ng := range ngrams2 {
		set2[ng] = true
	}

	// Find matching ngrams
	var matchCount int
	var matches []Match
	matchedPositions := make(map[int]bool)

	for i, ng := range ngrams1 {
		if set2[ng] {
			matchCount++
			if !matchedPositions[i] {
				matches = append(matches, Match{
					StartIndex:  i,
					EndIndex:    i + d.nGramSize,
					MatchedText: ng,
					Similarity:  1.0,
				})
				matchedPositions[i] = true
			}
		}
	}

	// Calculate Jaccard similarity
	unionSize := len(ngrams1) + len(ngrams2) - matchCount
	if unionSize == 0 {
		return 0, nil, nil
	}

	similarity := float64(matchCount) / float64(unionSize) * 100

	return similarity, mergeMatches(matches), nil
}

// CheckAgainstCorpus checks text against multiple submissions
func (d *detector) CheckAgainstCorpus(ctx context.Context, text string, submissionIDs []uuid.UUID) (*ReportDetails, error) {
	startTime := time.Now()

	// In a real implementation, this would fetch texts from database
	// and compare against each. For now, return placeholder.
	details := &ReportDetails{
		Matches:      []Match{},
		Sources:      []Source{},
		AnalysisTime: time.Since(startTime).Milliseconds(),
	}

	return details, nil
}

// normalize prepares text for comparison
func normalize(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Remove punctuation and extra whitespace
	var builder strings.Builder
	prevSpace := false
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			builder.WriteRune(r)
			prevSpace = false
		} else if unicode.IsSpace(r) && !prevSpace {
			builder.WriteRune(' ')
			prevSpace = true
		}
	}

	return strings.TrimSpace(builder.String())
}

// generateNGrams creates n-grams from text
func generateNGrams(text string, n int) []string {
	words := strings.Fields(text)
	if len(words) < n {
		return []string{strings.Join(words, " ")}
	}

	ngrams := make([]string, 0, len(words)-n+1)
	for i := 0; i <= len(words)-n; i++ {
		ngram := strings.Join(words[i:i+n], " ")
		ngrams = append(ngrams, ngram)
	}

	return ngrams
}

// mergeMatches combines adjacent matches
func mergeMatches(matches []Match) []Match {
	if len(matches) <= 1 {
		return matches
	}

	merged := []Match{matches[0]}
	for i := 1; i < len(matches); i++ {
		last := &merged[len(merged)-1]
		current := matches[i]

		// If adjacent or overlapping, merge
		if current.StartIndex <= last.EndIndex+1 {
			if current.EndIndex > last.EndIndex {
				last.EndIndex = current.EndIndex
				last.MatchedText += " " + current.MatchedText
			}
		} else {
			merged = append(merged, current)
		}
	}

	return merged
}
