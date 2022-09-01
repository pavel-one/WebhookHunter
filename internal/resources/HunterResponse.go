package resources

import "github.com/pavel-one/WebhookWatcher/internal/models"

type HunterResponse struct {
	URI string `json:"uri"`
}

func (r *HunterResponse) Init(hunter *models.Hunter) {
	r.URI = hunter.Slug
}
