package webscraper

import (
	"regexp"
	"strings"
)

// PreprocessMarkdownWebpage extracts core content from Markdown by removing navigation,
// menus, headers, footers while preserving the main content
func (s *webscraperService) processMarkdownWebpage(markdownContent string) string {
	// Phase 1: Remove navigation sections
	cleanText := s.removeNavigationSections(markdownContent)

	// Phase 2: Remove individual navigation elements
	cleanText = s.removeNavigationElements(cleanText)

	// Phase 3: Clean up and format the result
	return s.formatOutput(cleanText)
}

// removeNavigationSections removes entire sections that appear to be navigation
func (s *webscraperService) removeNavigationSections(text string) string {
	lines := strings.Split(text, "\n")
	var cleanLines []string
	inNavSection := false
	linkCounter := 0
	textCounter := 0

	// First pass: mark navigation-heavy sections
	for i, line := range lines {
		// Reset counters on headings and horizontal rules
		if s.isHeading(line) || s.isHorizontalRule(line) {
			// If we were in a nav section, decide whether to keep it
			if inNavSection && i > 0 {
				// If the section had more links than text, it was likely navigation
				if linkCounter > textCounter {
					// Remove the last few lines (navigation section)
					cleanLength := len(cleanLines)
					startTrim := max(0, cleanLength-linkCounter-1) // +1 for heading
					cleanLines = cleanLines[:startTrim]
				}
			}

			// Start fresh counting for the next section
			inNavSection = true
			linkCounter = 0
			textCounter = 0
		}

		// Count links vs text content
		if s.isMarkdownLink(line) || s.isHtmlLink(line) || s.isLinkDefinition(line) {
			linkCounter++
		} else if len(strings.TrimSpace(line)) > 0 && !s.isLinkListItem(line) {
			textCounter++
		}

		cleanLines = append(cleanLines, line)
	}

	// Second pass: filter out navigation markers and list-only sections
	filteredLines := []string{}
	inListSection := false
	listItemCount := 0
	textInSection := 0

	for i, line := range cleanLines {
		if s.isHeading(line) || s.isHorizontalRule(line) || i == len(cleanLines)-1 {
			// Evaluate the previous section
			if inListSection && listItemCount > 0 {
				// If this was just a list of links with no real content
				if textInSection == 0 || (listItemCount > textInSection*3) {
					// Go back and remove the list section
					newLen := len(filteredLines)
					cutPoint := newLen - listItemCount - 1 // +1 for heading
					if cutPoint < 0 {
						cutPoint = 0
					}
					filteredLines = filteredLines[:cutPoint]
				}
			}

			// Reset for new section
			inListSection = true
			listItemCount = 0
			textInSection = 0

			// Keep headings and horizontal rules
			filteredLines = append(filteredLines, line)
			continue
		}

		// Check for list items and track them
		if s.isListItem(line) {
			// Only count it as a list item if it looks like a navigation item
			if s.isLinkListItem(line) {
				listItemCount++
				continue // Skip adding it
			}
		} else if len(strings.TrimSpace(line)) > 0 {
			// Count non-empty, non-list lines as content
			textInSection++
		}

		// Keep lines that aren't nav-specific
		if !s.isNavigationMarker(line) && !s.isLinkSection(line) {
			filteredLines = append(filteredLines, line)
		}
	}

	return strings.Join(filteredLines, "\n")
}

// Helper functions to identify different elements

func (s *webscraperService) isHeading(line string) bool {
	headingRegex := regexp.MustCompile(`^#{1,6}\s+.+$|^.+\n[=\-]+$`)
	return headingRegex.MatchString(line)
}

func (s *webscraperService) isHorizontalRule(line string) bool {
	hrRegex := regexp.MustCompile(`^(\*{3,}|-{3,}|_{3,})$`)
	return hrRegex.MatchString(strings.TrimSpace(line))
}

func (s *webscraperService) isMarkdownLink(line string) bool {
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	return linkRegex.MatchString(line)
}

func (s *webscraperService) isHtmlLink(line string) bool {
	htmlLinkRegex := regexp.MustCompile(`<a\s+[^>]*>[^<]*<\/a>`)
	return htmlLinkRegex.MatchString(line)
}

func (s *webscraperService) isLinkDefinition(line string) bool {
	linkDefRegex := regexp.MustCompile(`^\s*\[[^\]]+\]:\s*http.+$`)
	return linkDefRegex.MatchString(line)
}

func (s *webscraperService) isListItem(line string) bool {
	listItemRegex := regexp.MustCompile(`^\s*[\*\-+]\s+.+$|^\s*\d+\.\s+.+$`)
	return listItemRegex.MatchString(line)
}

func (s *webscraperService) isLinkListItem(line string) bool {
	linkListItemRegex := regexp.MustCompile(`^\s*[\*\-+]\s+\[.+\]\(.+\).*$|^\s*\d+\.\s+\[.+\]\(.+\).*$`)
	return linkListItemRegex.MatchString(line)
}

func (s *webscraperService) isNavigationMarker(line string) bool {
	navMarkers := []string{
		"Navigation", "Menu", "Links", "Buttons", "Toggle navigation",
		"Copyright", "All rights reserved", "Follow us on", "Find us on",
		"Links/Buttons:", "Social media", "Back to", "Skip to",
	}

	lowered := strings.ToLower(line)
	for _, marker := range navMarkers {
		if strings.Contains(lowered, strings.ToLower(marker)) {
			return true
		}
	}
	return false
}

func (s *webscraperService) isLinkSection(line string) bool {
	linkSectionRegex := regexp.MustCompile(`^(Links|Nav|Menu|Navigation|Footer)(\W|$)`)
	return linkSectionRegex.MatchString(strings.TrimSpace(line))
}

// removeNavigationElements applies regex filters to remove specific navigation elements
func (s *webscraperService) removeNavigationElements(text string) string {
	// Patterns to identify and remove
	patterns := []*regexp.Regexp{
		// Links - only remove if they're a list of links with no surrounding context
		regexp.MustCompile(`(?m)^\s*\* \[[^\]]+\]\([^)]+\)\s*$`), // List item links
		regexp.MustCompile(`(?m)^\s*- \[[^\]]+\]\([^)]+\)\s*$`),  // List item links with dash
		regexp.MustCompile(`(?m)^\s*\+ \[[^\]]+\]\([^)]+\)\s*$`), // List item links with plus

		// Navigation lines and sections
		regexp.MustCompile(`(?i)toggle navigation`),           // Bootstrap toggle nav
		regexp.MustCompile(`(?i)copyright ©\d{4}.*`),          // Copyright notices
		regexp.MustCompile(`(?i)all rights reserved.*`),       // Rights reserved
		regexp.MustCompile(`(?i)(follow us on|find us on).*`), // Social media sections
		regexp.MustCompile(`(?m)^Links/Buttons:$.*?^$`),       // Links section headers
	}

	// Apply each pattern
	for _, pattern := range patterns {
		text = pattern.ReplaceAllString(text, "")
	}

	return text
}

// formatOutput cleans up the extracted text for better readability
func (s *webscraperService) formatOutput(text string) string {
	// Clean up excessive blank lines (more than 2 consecutive)
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	// Fix spacing around punctuation
	text = regexp.MustCompile(`\s+([.,;:!?])`).ReplaceAllString(text, "$1")

	// Keep important headings (= and - underlined headings)
	text = regexp.MustCompile(`(?m)^([^\n]+)\n[=]+\s*$`).ReplaceAllString(text, "$1\n")
	text = regexp.MustCompile(`(?m)^([^\n]+)\n[-]+\s*$`).ReplaceAllString(text, "$1\n")

	// Trim leading/trailing whitespace
	return strings.TrimSpace(text)
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
