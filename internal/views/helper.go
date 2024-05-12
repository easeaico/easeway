package views

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/easeaico/easeway/internal/config"
)

func RenderURL(path string) string {
	return fmt.Sprintf("%s://%s/%s", config.Conf.Site.Scheme, config.Conf.Site.Host, path)
}

func RenderSafeURL(path string) templ.SafeURL {
	return templ.URL(RenderURL(path))
}
