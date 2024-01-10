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
				Logo:               "/images/icon-32x32.svg",
				LogoAlt:            "DXTA",
				Content:            "DXTA has Economy growth upswing market index funds capitalization corporate bonds mutual.",
				BackgroundImage:    "https://images.unsplash.com/photo-1642132652859-3ef5a1048fd1?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2660&q=80",
				BackgroundImageAlt: "DXTA",
				URL:                "/",
			},
			SecondNarrowCard: &templates.Card{
				Logo:               "/images/icon-32x32.svg",
				LogoAlt:            "DXTA",
				Content:            "DXTA has Economy growth upswing market index funds capitalization corporate bonds mutual.",
				BackgroundImage:    "https://images.unsplash.com/photo-1522932753915-9ee97e43e3d9?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1yZWxhdGVkfDE0fHx8ZW58MHx8fHw%3D&auto=format&fit=crop&w=800&q=60",
				BackgroundImageAlt: "DXTA",
				URL:                "/",
			},
			WideCard: &templates.Card{
				Logo:               "/images/icon-32x32.svg",
				LogoAlt:            "DXTA",
				Content:            "DXTA has Economy growth upswing market index funds capitalization corporate bonds mutual.",
				BackgroundImage:    "https://images.unsplash.com/photo-1491895200222-0fc4a4c35e18?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1674&q=80",
				BackgroundImageAlt: "DXTA",
				URL:                "/",
			},
		},
		{
			FirstNarrowCard: &templates.Card{
				Logo:               "/images/icon-32x32.svg",
				LogoAlt:            "DXTA",
				Content:            "DXTA has Economy growth upswing market index funds capitalization corporate bonds mutual.",
				BackgroundImage:    "https://images.unsplash.com/photo-1642132652859-3ef5a1048fd1?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2660&q=80",
				BackgroundImageAlt: "DXTA",
				URL:                "/",
			},
			SecondNarrowCard: &templates.Card{
				Logo:               "/images/icon-32x32.svg",
				LogoAlt:            "DXTA",
				Content:            "DXTA has Economy growth upswing market index funds capitalization corporate bonds mutual.",
				BackgroundImage:    "https://images.unsplash.com/photo-1522932753915-9ee97e43e3d9?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1yZWxhdGVkfDE0fHx8ZW58MHx8fHw%3D&auto=format&fit=crop&w=800&q=60",
				BackgroundImageAlt: "DXTA",
				URL:                "/",
			},
			WideCard: &templates.Card{
				Logo:               "/images/icon-32x32.svg",
				LogoAlt:            "DXTA",
				Content:            "DXTA has Economy growth upswing market index funds capitalization corporate bonds mutual.",
				BackgroundImage:    "https://images.unsplash.com/photo-1491895200222-0fc4a4c35e18?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1674&q=80",
				BackgroundImageAlt: "DXTA",
				URL:                "/",
			},
		},
		{
			WideCard: &templates.Card{
				Logo:               "/images/icon-32x32.svg",
				LogoAlt:            "DXTA",
				Content:            "DXTA has Economy growth upswing market index funds capitalization corporate bonds mutual.",
				BackgroundImage:    "https://images.unsplash.com/photo-1642132652859-3ef5a1048fd1?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2660&q=80",
				BackgroundImageAlt: "DXTA",
				URL:                "/",
			},
		},
	}

	components := templates.OSSIndex(page, cardGroups)
	return components.Render(context.Background(), c.Response().Writer)
}
