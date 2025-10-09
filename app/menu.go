package app

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

func newMenu(menuItems [][]string) *tview.TextView {
	menu := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter)

	menu.SetBackgroundColor(style.BgColor)

	menuList := []string{}

	for i := range menuItems {
		key, item := genMenuItem(menuItems[i])
		if i == len(menuItems)-1 {
			item += " "
		}

		menuList = append(menuList, key+item)
	}

	_, err := fmt.Fprintf(menu, "%s", strings.Join(menuList, " "))
	if err != nil {
		log.Warn().Msgf("failed to create new menu: %s", err.Error())
	}

	return menu
}

func genMenuItem(items []string) (string, string) {
	key := fmt.Sprintf("[%s::b] <%s>[-:-:-]", style.GetColorHex(style.PageHeaderFgColor), items[0])
	// Translate and convert to uppercase
	translated := i18n.T(items[1])
	desc := fmt.Sprintf("[%s:%s:b] %s [-:-:-]",
		style.GetColorHex(style.PageHeaderFgColor),
		style.GetColorHex(style.MenuBgColor),
		strings.ToUpper(translated))

	return key, desc
}
