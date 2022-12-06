package models

import (
	"errors"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/pavel-one/WebhookWatcher/internal/helpers"
	"github.com/pavel-one/WebhookWatcher/internal/sqlite"
	"os"
)

type Hunter struct {
	Slug string `json:"slug"`
}

func (h *Hunter) Create() error {
	h.Slug = petname.Generate(2, "-") + "-" + helpers.RandString(5)
	_, err := os.Stat("./storage/users/")
	if err != nil {
		h.Slug = "default"
	}

	return h.CreateWithName()
}

func (h *Hunter) CreateWithName() (err error) {
	if h.Slug == "" {
		return errors.New("not set error")
	}

	path := "./storage/users/" + h.Slug
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	db, err := sqlite.GetDb(h.Slug)
	if err != nil {
		return err
	}

	err = sqlite.SetDefaultSchema(db)
	if err != nil {
		return err
	}

	err, _ = h.CreateChannel("/")
	if err != nil {
		return err
	}

	return nil
}

func (h *Hunter) FindBySlug(slug string) error {
	_, err := os.Stat("storage/users/" + slug)
	if err != nil {
		//if slug == "localhost" {
		//	h.Slug = "localhost"
		//	err := h.CreateWithName()
		//	if err != nil {
		//		return err
		//	}
		//} else {
		//	return err
		//}
		return err
	}

	h.Slug = slug
	return nil
}

func (h *Hunter) CreateChannel(channel string) (error, Channel) {
	if h.Slug == "" {
		return errors.New("hunter not exists"), Channel{}
	}

	db, err := sqlite.GetDb(h.Slug)
	if err != nil {
		return err, Channel{}
	}

	ch := Channel{
		Path: channel,
	}

	if err := ch.Create(db); err != nil {
		return err, Channel{}
	}

	return nil, ch
}

func (h *Hunter) FindChannelByPath(channel string) (error, Channel) {
	var ch Channel

	db, err := sqlite.GetDb(h.Slug)
	if err != nil {
		return err, Channel{}
	}

	err = db.Get(&ch, "SELECT * FROM channels WHERE path=$1 ORDER BY created_at desc", channel)
	if err != nil {
		return err, ch
	}

	return nil, ch
}

func (h *Hunter) Channels() ([]Channel, error) {
	var channels []Channel

	db, err := sqlite.GetDb(h.Slug)
	if err != nil {
		return nil, err
	}

	err = db.Select(&channels, `select c.id, c.path, c.created_at, count(r.id) request_count
										from channels as c
												 left join requests r on c.id = r.channel_id
										group by c.id
										order by c.created_at desc`)
	if err != nil {
		return nil, err
	}

	return channels, nil
}
