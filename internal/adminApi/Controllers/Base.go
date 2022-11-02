package AdminControllers

import (
	"errors"
	"github.com/pavel-one/WebhookWatcher/internal/controllers"
	"github.com/pavel-one/WebhookWatcher/internal/models"
	"net/http"
)

type SubController struct {
	controllers.BaseController
	controllers.DatabaseController
}

func (c *SubController) checkHunter(slug string) (*models.Hunter, error) {
	hunter := new(models.Hunter)
	hunter.FindBySlug(c.DB, slug)

	if hunter.Id == "" {
		return nil, errors.New(hunterErr)
	}

	return hunter, nil
}

func (c *SubController) checkChannel(hunter *models.Hunter, path string) (models.Channel, int, error) {
	channel, err := hunter.FindChannelByPath(c.DB, "/"+path)

	if err != nil {
		return channel, http.StatusNotFound, err
	}

	if channel.Id == 0 {
		return channel, http.StatusNotFound, errors.New("channel not found")
	}

	return channel, http.StatusOK, nil
}
