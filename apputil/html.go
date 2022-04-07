package apputil

import (
	"bytes"
	"io"
	"io/fs"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"tractor.dev/apptron/client"
)

// OptionsFromHTML will use fsys to open filename and parse it as HTML to find a meta element within
// the head element with metaname for its name attribute and produce a WindowOptions populated with
// the comma-separated key-values of the meta elements content attribute. If there is no key-value for
// a field of WindowOptions, or is the HTML, Icon, or Script field, the value from fallback is used.
//
// If a title element is found within the head element, its value will be used. If the meta element comes
// after the title element and includes a key-value for title, it will be used instead.
//
// All keys in the content attribute are lowercase equivalents of their field name. However, for Position.X
// and Position.Y, keys x and y are used. Size.Width and Size.Height keys are width and height. Keys for
// MinSize and MaxSize fields are min-width, min-height, max-width, and max-height, respectively.
//
// String values can optionally be wrapped in single quotes, boolean values must be true or false. Any errors
// in parsing values will result in using the fallback struct value.
func OptionsFromHTML(fsys fs.FS, filename, metaname string, fallback client.WindowOptions) client.WindowOptions {
	f, err := fsys.Open(filename)
	if err != nil {
		return fallback
	}
	opts := parseOptions(f, metaname)
	f.Close()
	return client.WindowOptions{
		Title:       optString(opts["title"], fallback.Title),
		AlwaysOnTop: optBool(opts["alwaysontop"], fallback.AlwaysOnTop),
		Fullscreen:  optBool(opts["fullscreen"], fallback.Fullscreen),
		Maximized:   optBool(opts["maximized"], fallback.Maximized),
		Resizable:   optBool(opts["resizable"], fallback.Resizable),
		Transparent: optBool(opts["transparent"], fallback.Transparent),
		Frameless:   optBool(opts["frameless"], fallback.Frameless),
		Visible:     optBool(opts["visible"], fallback.Visible),
		Center:      optBool(opts["center"], fallback.Center),
		Position: client.Position{
			X: optFloat(opts["x"], fallback.Position.X),
			Y: optFloat(opts["y"], fallback.Position.Y),
		},
		Size: client.Size{
			Width:  optFloat(opts["width"], fallback.Size.Width),
			Height: optFloat(opts["height"], fallback.Size.Height),
		},
		MinSize: client.Size{
			Width:  optFloat(opts["min-width"], fallback.MinSize.Width),
			Height: optFloat(opts["min-height"], fallback.MinSize.Height),
		},
		MaxSize: client.Size{
			Width:  optFloat(opts["max-width"], fallback.MaxSize.Width),
			Height: optFloat(opts["max-height"], fallback.MaxSize.Height),
		},
		URL:     optString(opts["url"], fallback.URL),
		Script:  fallback.Script,
		HTML:    fallback.HTML,
		IconSel: fallback.IconSel,
		Icon:    fallback.Icon,
	}
}

func parseOptions(f io.Reader, metaname string) map[string]string {
	z := html.NewTokenizer(f)
	inTitle := false
	opts := make(map[string]string)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// also EOF aka done
			return opts
		case html.StartTagToken:
			tn, _ := z.TagName()
			if bytes.Equal(tn, []byte("body")) {
				return opts
			}
			if bytes.Equal(tn, []byte("title")) {
				inTitle = true
			}
		case html.TextToken:
			if inTitle {
				opts["title"] = string(z.Text())
				inTitle = false
			}
		case html.SelfClosingTagToken:
			tn, attr := z.TagName()
			if bytes.Equal(tn, []byte("meta")) && attr {
				var name, content string
				for attr {
					var k, v []byte
					k, v, attr = z.TagAttr()
					if bytes.Equal(k, []byte("name")) {
						name = string(v)
					}
					if bytes.Equal(k, []byte("content")) {
						content = string(v)
					}
				}
				if name == metaname {
					for _, part := range strings.Split(content, ",") {
						kv := strings.SplitN(part, "=", 2)
						opts[kv[0]] = kv[1]
					}
				}
			}
		}
	}
}

func optBool(v string, fallback bool) bool {
	if v == "true" {
		return true
	}
	if v == "false" {
		return false
	}
	return fallback
}

func optFloat(v string, fallback float64) float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return f
	}
	return fallback
}

func optString(v string, fallback string) string {
	if v == "" {
		return fallback
	}
	return strings.Trim(v, "'")
}
