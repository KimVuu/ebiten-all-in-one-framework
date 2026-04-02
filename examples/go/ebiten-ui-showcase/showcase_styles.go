package main

import ebitenui "github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui"

func showcaseGroupStyle() ebitenui.Style {
	return showcaseGroupStyleForChrome(buildDefaultShowcasePreset().Chrome)
}

func showcaseGroupStyleForChrome(chrome showcaseChrome) ebitenui.Style {
	return ebitenui.Style{
		Width:           ebitenui.Fill(),
		Direction:       ebitenui.Column,
		Padding:         ebitenui.All(16),
		Gap:             12,
		BackgroundColor: chrome.PanelBackground,
		BorderColor:     chrome.PanelBorder,
		BorderWidth:     1,
	}
}

func showcaseGroupTitleStyle() ebitenui.Style {
	return showcaseGroupTitleStyleForChrome(buildDefaultShowcasePreset().Chrome)
}

func showcaseGroupTitleStyleForChrome(chrome showcaseChrome) ebitenui.Style {
	return ebitenui.Style{
		Color: chrome.TextStrong,
	}
}

func showcaseGroupCopyStyle() ebitenui.Style {
	return showcaseGroupCopyStyleForChrome(buildDefaultShowcasePreset().Chrome)
}

func showcaseGroupCopyStyleForChrome(chrome showcaseChrome) ebitenui.Style {
	return ebitenui.Style{
		Width:      ebitenui.Fill(),
		Color:      chrome.TextMuted,
		LineHeight: 16,
	}
}

func detailSectionStyle() ebitenui.Style {
	return detailSectionStyleForChrome(buildDefaultShowcasePreset().Chrome)
}

func detailSectionStyleForChrome(chrome showcaseChrome) ebitenui.Style {
	return ebitenui.Style{
		Width:           ebitenui.Fill(),
		Direction:       ebitenui.Column,
		Padding:         ebitenui.All(16),
		Gap:             12,
		BackgroundColor: chrome.PanelBackground,
		BorderColor:     chrome.PanelBorder,
		BorderWidth:     1,
	}
}

func detailTitleStyle() ebitenui.Style {
	return detailTitleStyleForChrome(buildDefaultShowcasePreset().Chrome)
}

func detailTitleStyleForChrome(chrome showcaseChrome) ebitenui.Style {
	return ebitenui.Style{
		Color: chrome.TextStrong,
	}
}
