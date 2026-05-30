package media

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BerrySDK/berryone/events"
)

type LoadedMedia struct {
	Buffer   []byte
	FileName string
	Mimetype string
}

type Manager struct {
	HTTPClient *http.Client
}

func NewManager() *Manager {
	return &Manager{
		HTTPClient: http.DefaultClient,
	}
}

func (m *Manager) Load(payload events.MediaPayload) (LoadedMedia, error) {
	switch {
	case len(payload.Buffer) > 0:
		return LoadedMedia{
			Buffer:   payload.Buffer,
			FileName: payload.FileName,
			Mimetype: payload.Mimetype,
		}, nil
	case payload.Path != "":
		buffer, err := os.ReadFile(payload.Path)
		if err != nil {
			return LoadedMedia{}, err
		}
		fileName := payload.FileName
		if fileName == "" {
			fileName = filepath.Base(payload.Path)
		}
		return LoadedMedia{
			Buffer:   buffer,
			FileName: fileName,
			Mimetype: payload.Mimetype,
		}, nil
	case payload.URL != "":
		response, err := m.HTTPClient.Get(payload.URL)
		if err != nil {
			return LoadedMedia{}, err
		}
		defer response.Body.Close()
		buffer, err := io.ReadAll(response.Body)
		if err != nil {
			return LoadedMedia{}, err
		}
		fileName := payload.FileName
		if fileName == "" {
			fileName = filepath.Base(response.Request.URL.Path)
		}
		return LoadedMedia{
			Buffer:   buffer,
			FileName: fileName,
			Mimetype: firstNonEmpty(payload.Mimetype, response.Header.Get("Content-Type")),
		}, nil
	default:
		return LoadedMedia{}, errors.New("media payload must provide buffer, path or url")
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
