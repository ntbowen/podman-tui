package help

import (
	"fmt"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Help is a help primitive dialog.
type Help struct {
	*tview.Box

	title    string
	layout   *tview.Flex
	keyinfo  *tview.Table
	appinfo  *tview.TextView
	appName  string
	appVer   string
}

// NewHelp returns a help screen primitive.
func NewHelp(appName string, appVersion string) *Help {
	// returns the help primitive
	help := &Help{
		Box:     tview.NewBox(),
		title:   "help",
		appName: appName,
		appVer:  appVersion,
	}

	// colors
	headerColor := style.HelpHeaderFgColor
	fgColor := style.FgColor
	bgColor := style.BgColor
	borderColor := style.BorderColor

	// application keys description table
	help.keyinfo = tview.NewTable()
	help.keyinfo.SetBackgroundColor(bgColor)
	help.keyinfo.SetFixed(1, 1)
	help.keyinfo.SetSelectable(false, false)

	// application description and version text view
	help.appinfo = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	help.appinfo.SetBackgroundColor(bgColor)

	licenseInfo := i18n.T("released under the Apache License 2.0.")
	appInfoText := fmt.Sprintf("%s %s - %s", appName, appVersion, licenseInfo)

	help.appinfo.SetText(appInfoText)
	help.appinfo.SetTextColor(headerColor)

	// help table items
	// the items will be divided into two separate tables
	rowIndex := 0
	colIndex := 0
	needInit := true
	maxRowIndex := len(utils.UIKeysBindings) / 2 //nolint:mnd

	for i := range utils.UIKeysBindings {
		if i >= maxRowIndex {
			if needInit {
				colIndex = 2
				rowIndex = 0
				needInit = false
			}
		}

		help.keyinfo.SetCell(rowIndex, colIndex,
			tview.NewTableCell(fmt.Sprintf("%s:", utils.UIKeysBindings[i].KeyLabel)). //nolint:perfsprint
													SetAlign(tview.AlignRight).
													SetBackgroundColor(bgColor).
													SetSelectable(true).SetTextColor(headerColor))

		help.keyinfo.SetCell(rowIndex, colIndex+1,
			tview.NewTableCell(i18n.T(utils.UIKeysBindings[i].KeyDesc)).
				SetAlign(tview.AlignLeft).
				SetBackgroundColor(bgColor).
				SetSelectable(true).SetTextColor(fgColor))

		rowIndex++
	}

	// appinfo and appkeys layout
	mlayout := tview.NewFlex().SetDirection(tview.FlexRow)
	mlayout.AddItem(help.appinfo, 1, 0, false)
	mlayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	mlayout.AddItem(help.keyinfo, 0, 1, false)
	mlayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	// layout
	help.layout = tview.NewFlex().SetDirection(tview.FlexColumn)
	help.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	help.layout.AddItem(mlayout, 0, 1, false)
	help.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	help.layout.SetBorder(true)
	help.layout.SetBackgroundColor(bgColor)
	help.layout.SetBorderColor(borderColor)

	return help
}

// GetTitle returns primitive title.
func (help *Help) GetTitle() string {
	return help.title
}

// HasFocus returns whether or not this primitive has focus.
func (help *Help) HasFocus() bool {
	return help.Box.HasFocus() || help.layout.HasFocus()
}

// Focus is called when this primitive receives focus.
func (help *Help) Focus(delegate func(p tview.Primitive)) {
	delegate(help.layout)
}

// RefreshLabels refreshes all key descriptions with current language translations.
func (help *Help) RefreshLabels() {
	// Update license info
	licenseInfo := i18n.T("released under the Apache License 2.0.")
	appInfoText := fmt.Sprintf("%s %s - %s", help.appName, help.appVer, licenseInfo)
	help.appinfo.SetText(appInfoText)
	
	// Update key descriptions
	rowIndex := 0
	colIndex := 0
	needInit := true
	maxRowIndex := len(utils.UIKeysBindings) / 2 //nolint:mnd

	for i := range utils.UIKeysBindings {
		if i >= maxRowIndex {
			if needInit {
				colIndex = 2
				rowIndex = 0
				needInit = false
			}
		}

		// Update only the description cell with translated text
		help.keyinfo.GetCell(rowIndex, colIndex+1).SetText(i18n.T(utils.UIKeysBindings[i].KeyDesc))
		rowIndex++
	}
}

// Draw draws this primitive onto the screen.
func (help *Help) Draw(screen tcell.Screen) {
	x, y, width, height := help.GetInnerRect()
	if height <= 3 { //nolint:mnd
		return
	}

	help.DrawForSubclass(screen, help)
	help.layout.SetRect(x, y, width, height)
	help.layout.Draw(screen)
}
