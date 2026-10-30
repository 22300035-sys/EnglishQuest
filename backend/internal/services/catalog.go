package services

import (
	"encoding/json"
	"github.com/englishquest/backend/curriculum"
	"github.com/englishquest/backend/internal/models"
	apperr "github.com/englishquest/backend/internal/pkg/errors"
	"net"
	"net/url"
	"regexp"
	"strings"
)

var ResourceKey = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,39}$`)
var localImage = regexp.MustCompile(`^/images/[a-zA-Z0-9_-]+\.(svg|png|jpg|webp)$`)
var localAudio = regexp.MustCompile(`^/audio/[a-zA-Z0-9_-]+\.(wav|mp3|ogg)$`)
var BundledAudio = map[string]bool{
	"/audio/eq-listen-a1-1.wav": true, "/audio/eq-listen-a1-2.wav": true, "/audio/eq-listen-a1-3.wav": true,
	"/audio/eq-listen-a2-1.wav": true, "/audio/eq-listen-a2-2.wav": true, "/audio/eq-listen-a2-3.wav": true,
}

func init() {
	var units []struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(curriculum.Data, &units); err != nil {
		panic("invalid bundled curriculum: " + err.Error())
	}
	for _, unit := range units {
		BundledAudio["/audio/course-"+unit.Key+".wav"] = true
	}
}

func ValidateResource(in models.Resource) error {
	bad := func() error {
		return apperr.NewAppError(400, "Recurso inválido. Usa una clave simple, un título, atribución y una URL HTTPS pública o un archivo local permitido.")
	}
	if !ResourceKey.MatchString(in.Key) || strings.TrimSpace(in.Title) == "" || len(in.Title) > 150 || len(in.Attribution) < 3 || len(in.Attribution) > 2000 || len(in.URL) > 500 {
		return bad()
	}
	if in.Kind != "image" && in.Kind != "audio" {
		return bad()
	}
	if strings.HasPrefix(in.URL, "/") {
		if (in.Kind == "image" && !localImage.MatchString(in.URL)) || (in.Kind == "audio" && !localAudio.MatchString(in.URL)) {
			return bad()
		}
		if in.Kind == "image" && (!strings.HasSuffix(in.URL, ".svg") || !Images[strings.TrimSuffix(strings.TrimPrefix(in.URL, "/images/"), ".svg")]) {
			return bad()
		}
		if in.Kind == "audio" && !BundledAudio[in.URL] {
			return bad()
		}
		return nil
	}
	u, err := url.Parse(in.URL)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Hostname() == "localhost" {
		return bad()
	}
	if ip := net.ParseIP(u.Hostname()); ip != nil && (ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast()) {
		return bad()
	}
	return nil
}
