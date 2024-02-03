package handlers

import (
	"context"
	"dxta-dev/app/internal/templates"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) OSSIndex(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "oss",
		Boosted: h.HxBoosted,
	}

	cardGroups := []templates.CardGroup{
		{
			FirstNarrowCard: &templates.Card{
				Logo:               "/images/oss/pocketbase-logo.svg",
				LogoAlt:            "Pocketbase",
				Content:            "Pocketbase is open source backend for your next SaaS and Mobile app in 1 file",
				BackgroundImage:    "/images/oss/pocketbase-background.png",
				BackgroundImageAlt: "Pocketbase",
				URL:                "https://pocketbase.dxta.dev",
			},
			SecondNarrowCard: nil,
			WideCard: &templates.Card{
				Logo:               "/images/oss/calcom-logo.svg",
				LogoAlt:            "DXTA",
				Content:            "Cal.com is the event-juggling scheduler for everyone. Focus on meeting, not making meetings.",
				BackgroundImage:    "/images/oss/calcom-background.png",
				BackgroundImageAlt: "DXTA",
				URL:                "https://calcom.dxta.dev",
			},
		},
		{
			FirstNarrowCard: &templates.Card{
				Logo:               "/images/oss/documenso-logo.svg",
				LogoAlt:            "Documenso",
				Content:            "Documenso is an open sourced document signing platform that allows you to sign documents with ease.",
				BackgroundImage:    "/images/oss/documenso-background.png",
				BackgroundImageAlt: "DXTA",
				URL:                "https://documenso.dxta.dev",
			},
			SecondNarrowCard: &templates.Card{
				Logo:               "/images/oss/dub-logo.svg",
				LogoAlt:            "DUB",
				Content:            "DUB is a simple, open source, and free to use URL shortener.",
				BackgroundImage:    "/images/oss/dub-background.png",
				BackgroundImageAlt: "DUB",
				URL:                "https://dub.dxta.dev",
			},
			WideCard: nil,
		},
		{
			WideCard: nil,
		},
	}

	components := templates.OSSIndex(page, cardGroups)
	return components.Render(context.Background(), c.Response().Writer)
}
