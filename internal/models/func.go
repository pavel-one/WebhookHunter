package models

import "github.com/dustin/go-humanize"

func FormatRequests(requests []RequestModel) []map[string]any {
	var out []map[string]any

	for _, item := range requests {
		out = append(out, FormatRequest(item))
	}

	return out
}

func FormatRequest(item RequestModel) map[string]any {
	return map[string]any{
		"id":         item.ID,
		"request":    item.Request,
		"headers":    item.Headers,
		"query":      item.Query,
		"created_at": item.CreatedAt.Format("15:04:05 02.01.2006"),
		"humanize":   humanize.Time(item.CreatedAt),
	}
}
