package images

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/i18n"
	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/images/imgdialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

const (
	viewImageRepoNameColIndex = 0 + iota
	viewImageTagColIndex
	viewImageIDColIndex
	viewImageCreatedAtColIndex
	viewImageSizeColIndex
)

var (
	errNoImageToTree       = errors.New("here is no image to display tree")
	errNoImageToUntag      = errors.New("here is no image to untag")
	errNoImageToTag        = errors.New("here is no image to tag")
	errNoImageToSave       = errors.New("here is no image to save")
	errNoImageToPush       = errors.New("here is no image to push")
	errNoImageToHistory    = errors.New("here is no image to display history")
	errNoImageToDiff       = errors.New("here is no image to display diff")
	errNoImageToRemove     = errors.New("there is no image to remove")
	errNoImageToInspect    = errors.New("there is no image to display inspect")
	errNoBuildDirOrCntFile = errors.New("both context directory path and container files fields are empty")
)

// Images implements the images primitive.
type Images struct {
	*tview.Box

	title           string
	headers         []string
	table           *tview.Table
	errorDialog     *dialogs.ErrorDialog
	progressDialog  *dialogs.ProgressDialog
	cmdDialog       *dialogs.CommandDialog
	cmdInputDialog  *dialogs.SimpleInputDialog
	messageDialog   *dialogs.MessageDialog
	confirmDialog   *dialogs.ConfirmDialog
	sortDialog      *dialogs.SortDialog
	searchDialog    *imgdialogs.ImageSearchDialog
	historyDialog   *imgdialogs.ImageHistoryDialog
	importDialog    *imgdialogs.ImageImportDialog
	buildDialog     *imgdialogs.ImageBuildDialog
	buildPrgDialog  *imgdialogs.ImageBuildProgressDialog
	saveDialog      *imgdialogs.ImageSaveDialog
	pushDialog      *imgdialogs.ImagePushDialog
	imagesList      imageListReport
	selectedID      string
	selectedName    string
	confirmData     string
	fastRefreshChan chan bool
	appFocusHandler func()
}

type imageListReport struct {
	mu        sync.Mutex
	report    []images.ImageListReporter
	sortBy    string
	ascending bool
}

// NewImages returns images page view.
func NewImages() *Images {
	images := &Images{
		Box:            tview.NewBox(),
		title:          "images",
		headers:        []string{i18n.T("repository"), i18n.T("tag"), i18n.T("image id"), i18n.T("created at"), i18n.T("size")},
		errorDialog:    dialogs.NewErrorDialog(),
		progressDialog: dialogs.NewProgressDialog(),
		cmdInputDialog: dialogs.NewSimpleInputDialog(""),
		messageDialog:  dialogs.NewMessageDialog(""),
		confirmDialog:  dialogs.NewConfirmDialog(),
		sortDialog:     dialogs.NewSortDialog([]string{i18n.T("repository"), i18n.T("created"), i18n.T("size")}, 1),
		searchDialog:   imgdialogs.NewImageSearchDialog(),
		historyDialog:  imgdialogs.NewImageHistoryDialog(),
		importDialog:   imgdialogs.NewImageImportDialog(),
		buildDialog:    imgdialogs.NewImageBuildDialog(),
		buildPrgDialog: imgdialogs.NewImageBuildProgressDialog(),
		saveDialog:     imgdialogs.NewImageSaveDialog(),
		pushDialog:     imgdialogs.NewImagePushDialog(),
		imagesList:     imageListReport{sortBy: "created", ascending: true},
	}

	// Build command dialog with translations
	images.buildCommandDialog()

	imgTable := tview.NewTable()
	translatedTitle := i18n.T(images.title)
	imgTable.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(translatedTitle)))
	imgTable.SetBorderColor(style.BorderColor)
	imgTable.SetBackgroundColor(style.BgColor)
	imgTable.SetTitleColor(style.FgColor)
	imgTable.SetBorder(true)

	for i := range images.headers {
		imgTable.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(images.headers[i]))). //nolint:perfsprint
														SetExpansion(1).
														SetBackgroundColor(style.PageHeaderBgColor).
														SetTextColor(style.PageHeaderFgColor).
														SetAlign(tview.AlignLeft).
														SetSelectable(false))
	}

	imgTable.SetFixed(1, 1)
	imgTable.SetSelectable(true, false)
	images.table = imgTable

	// set message dialog functions
	images.messageDialog.SetCancelFunc(func() {
		images.messageDialog.Hide()
	})

	// set input cmd dialog functions
	images.cmdInputDialog.SetCancelFunc(func() {
		images.cmdInputDialog.Hide()
	})

	images.cmdInputDialog.SetSelectedFunc(func() {
		images.cmdInputDialog.Hide()
	})

	// NOTE: cmdDialog handlers (SetSelectedFunc, SetCancelFunc) are set in buildCommandDialog()

	// set confirm dialogs functions
	images.confirmDialog.SetSelectedFunc(func() {
		images.confirmDialog.Hide()

		switch images.confirmData {
		case utils.PruneCommandLabel:
			images.prune()
		case "rm":
			images.remove()
		}
	})

	images.confirmDialog.SetCancelFunc(func() {
		images.confirmDialog.Hide()
	})

	// set history dialogs functions
	images.historyDialog.SetCancelFunc(func() {
		images.historyDialog.Hide()
	})

	// set search dialogs functions
	images.searchDialog.SetCancelFunc(func() {
		images.searchDialog.Hide()
	})

	images.searchDialog.SetSearchFunc(func() {
		term := images.searchDialog.GetSearchText()
		if term == "" {
			return
		}

		images.search(term)
	})

	images.searchDialog.SetPullFunc(func() {
		name := images.searchDialog.GetSelectedItem()
		images.pull(name)
	})

	// set build dialogs functions
	images.buildDialog.SetCancelFunc(images.buildDialog.Hide)
	images.buildDialog.SetBuildFunc(images.build)
	images.buildPrgDialog.SetFastRefreshHandler(func() {
		images.fastRefreshChan <- true
	})

	// set save dialog functions
	images.saveDialog.SetCancelFunc(images.saveDialog.Hide)
	images.saveDialog.SetSaveFunc(images.save)

	// set import dialog functions
	images.importDialog.SetCancelFunc(images.importDialog.Hide)
	images.importDialog.SetImportFunc(images.imageImport)

	// set push dialog functions
	images.pushDialog.SetPushFunc(images.push)
	images.pushDialog.SetCancelFunc(images.pushDialog.Hide)

	// set sort dialog functions
	images.sortDialog.SetSelectFunc(images.SortView)
	images.sortDialog.SetCancelFunc(images.sortDialog.Hide)

	return images
}

// SetAppFocusHandler sets application focus handler.
func (img *Images) SetAppFocusHandler(handler func()) {
	img.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (img *Images) GetTitle() string {
	return img.title
}

// HasFocus returns whether or not this primitive has focus.
func (img *Images) HasFocus() bool {
	if img.SubDialogHasFocus() {
		return true
	}

	if img.table.HasFocus() || img.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (img *Images) SubDialogHasFocus() bool {
	for _, dialog := range img.getInnerDialogs() {
		if dialog.HasFocus() {
			return true
		}
	}

	for _, dialog := range img.getInnerTopDialogs() {
		if dialog.HasFocus() {
			return true
		}
	}

	return false
}

// Focus is called when this primitive receives focus.
func (img *Images) Focus(delegate func(p tview.Primitive)) {
	// since error and confirm dialog can get focus on top of other dialogs
	if img.errorDialog.IsDisplay() {
		delegate(img.errorDialog)

		return
	}

	if img.confirmDialog.IsDisplay() {
		delegate(img.confirmDialog)

		return
	}

	for _, dialog := range img.getInnerDialogs() {
		if dialog.IsDisplay() {
			delegate(dialog)

			return
		}
	}

	for _, dialog := range img.getInnerTopDialogs() {
		if dialog.IsDisplay() {
			delegate(dialog)

			return
		}
	}

	delegate(img.table)
}

// HideAllDialogs hides all sub dialogs.
func (img *Images) HideAllDialogs() {
	for _, dialog := range img.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.Hide()
		}
	}

	for _, dialog := range img.getInnerTopDialogs() {
		if dialog.IsDisplay() {
			dialog.Hide()
		}
	}
}

// SetFastRefreshChannel sets channel for fastRefresh func.
func (img *Images) SetFastRefreshChannel(refresh chan bool) {
	img.fastRefreshChan = refresh
}

func (img *Images) getSelectedItem() (string, string) {
	if img.table.GetRowCount() <= 1 {
		return "", ""
	}

	row, _ := img.table.GetSelection()
	imageRepo := img.table.GetCell(row, 0).Text
	imageTag := img.table.GetCell(row, 1).Text
	imageName := imageRepo + ":" + imageTag
	imageID := img.table.GetCell(row, 2).Text //nolint:mnd

	return imageID, imageName
}

func (img *Images) getInnerDialogs() []utils.UIDialog {
	dialogs := []utils.UIDialog{
		img.cmdDialog,
		img.cmdInputDialog,
		img.messageDialog,
		img.searchDialog,
		img.historyDialog,
		img.importDialog,
		img.buildDialog,
		img.buildPrgDialog,
		img.saveDialog,
		img.pushDialog,
		img.sortDialog,
	}

	return dialogs
}

func (img *Images) getInnerTopDialogs() []utils.UIDialog {
	dialogs := []utils.UIDialog{
		img.errorDialog,
		img.progressDialog,
		img.confirmDialog,
	}

	return dialogs
}

// buildCommandDialog builds the command dialog with translated strings
func (img *Images) buildCommandDialog() {
	if img.cmdDialog != nil {
		img.cmdDialog.Hide()
	}

	// Translate both command names and descriptions
	img.cmdDialog = dialogs.NewCommandDialog([][]string{
		{i18n.T("build"), i18n.T("build an image from Containerfile")},
		{i18n.T("diff"), i18n.T("inspect changes to the image's file systems")},
		{i18n.T("history"), i18n.T("show history of the selected image")},
		{i18n.T("import"), i18n.T("create a container image from a tarball")},
		{i18n.T("inspect"), i18n.T("display the configuration of the selected image")},
		{i18n.T("prune"), i18n.T("remove all unused images")},
		{i18n.T("push"), i18n.T("push a source image to a specified destination")},
		{i18n.T("rm"), i18n.T("removes the selected  image from local storage")},
		{i18n.T("save"), i18n.T("save an image to docker-archive or oci-archive")},
		{i18n.T("search/pull"), i18n.T("search and pull image from registry")},
		{i18n.T("tag"), i18n.T("add an additional name to the selected  image")},
		{i18n.T("tree"), i18n.T("display layer hierarchy of an image")},
		{i18n.T("untag"), i18n.T("remove a name from the selected image")},
	})

	img.cmdDialog.SetSelectedFunc(func() {
		img.cmdDialog.Hide()
		// GetSelectedItem returns translated command, map it back to English
		translatedCmd := img.cmdDialog.GetSelectedItem()
		englishCmd := img.getEnglishCommand(translatedCmd)
		img.runCommand(englishCmd)
	})

	img.cmdDialog.SetCancelFunc(func() {
		img.cmdDialog.Hide()
	})
}

// getEnglishCommand maps translated command back to English command key
func (img *Images) getEnglishCommand(translatedCmd string) string {
	// Create reverse mapping from translated to English
	commandMap := map[string]string{
		i18n.T("build"):       "build",
		i18n.T("diff"):        "diff",
		i18n.T("history"):     "history",
		i18n.T("import"):      "import",
		i18n.T("inspect"):     "inspect",
		i18n.T("prune"):       "prune",
		i18n.T("push"):        "push",
		i18n.T("rm"):          "rm",
		i18n.T("save"):        "save",
		i18n.T("search/pull"): "search/pull",
		i18n.T("tag"):         "tag",
		i18n.T("tree"):        "tree",
		i18n.T("untag"):       "untag",
	}

	if englishCmd, exists := commandMap[translatedCmd]; exists {
		return englishCmd
	}

	// Fallback to original if not found (for English or unknown commands)
	return translatedCmd
}

// UpdateLanguage updates all translatable text when language changes
func (img *Images) UpdateLanguage() {
	// Update headers
	img.headers = []string{
		i18n.T("repository"),
		i18n.T("tag"),
		i18n.T("image id"),
		i18n.T("created at"),
		i18n.T("size"),
	}

	// Rebuild command dialog with new translations
	img.buildCommandDialog()

	// Update table header cells
	for i := range img.headers {
		header := fmt.Sprintf("[black::b]%s", strings.ToUpper(img.headers[i]))
		img.table.GetCell(0, i).SetText(header)
	}

	// Update table title with translation
	translatedTitle := i18n.T(img.title)
	img.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(translatedTitle), img.table.GetRowCount()-1))

	// Update dialogs
	img.messageDialog.UpdateLanguage()
	img.confirmDialog.UpdateLanguage()
	img.buildDialog.UpdateLanguage()
	img.buildPrgDialog.UpdateLanguage()
	img.historyDialog.UpdateLanguage()
	img.importDialog.UpdateLanguage()
	img.pushDialog.UpdateLanguage()
	img.saveDialog.UpdateLanguage()
	img.searchDialog.UpdateLanguage()

	// Update error dialog
	img.errorDialog.UpdateLanguage()
	
	// Rebuild sort dialog to update button text
	img.sortDialog = dialogs.NewSortDialog([]string{
		i18n.T("repository"),
		i18n.T("created"),
		i18n.T("size"),
	}, 1)

	// Re-set sort dialog handlers after rebuild
	img.sortDialog.SetSelectFunc(img.SortView)
	img.sortDialog.SetCancelFunc(img.sortDialog.Hide)
}
