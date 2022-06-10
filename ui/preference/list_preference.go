package preference

import (
	"sort"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type ListPreferenceModal struct {
	*load.Load
	*decredmaterial.Modal

	optionsRadioGroup *widget.Enum
	cancelButton      decredmaterial.IconButton

	items         map[string]string //[key]str-key
	itemKeys      []string
	title         string
	preferenceKey string
	defaultValue  string // str-key
	initialValue  string
	currentValue  string

	updateButtonClicked func()
}

func NewListPreference(l *load.Load, preferenceKey, defaultValue string, items map[string]string) *ListPreferenceModal {

	// sort keys to keep order when refreshed
	sortedKeys := make([]string, 0)
	for k := range items {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Slice(sortedKeys, func(i int, j int) bool { return sortedKeys[i] < sortedKeys[j] })

	lp := ListPreferenceModal{
		Load:          l,
		preferenceKey: preferenceKey,
		defaultValue:  defaultValue,

		items:    items,
		itemKeys: sortedKeys,

		optionsRadioGroup: new(widget.Enum),
		Modal:             l.Theme.ModalFloatTitle("list_preference"),
	}

	lp.cancelButton, _ = components.SubpageHeaderButtons(l)
	lp.cancelButton.Icon = l.Theme.Icons.ContentClear

	return &lp
}

func (lp *ListPreferenceModal) OnResume() {
	initialValue := lp.WL.MultiWallet.ReadStringConfigValueForKey(lp.preferenceKey)
	if initialValue == "" {
		initialValue = lp.defaultValue
	}

	lp.initialValue = initialValue
	lp.currentValue = initialValue

	lp.optionsRadioGroup.Value = lp.currentValue
}

func (lp *ListPreferenceModal) OnDismiss() {}

func (lp *ListPreferenceModal) Title(title string) *ListPreferenceModal {
	lp.title = title
	return lp
}

func (lp *ListPreferenceModal) UpdateValues(clicked func()) *ListPreferenceModal {
	lp.updateButtonClicked = clicked
	return lp
}

func (lp *ListPreferenceModal) Handle() {

	for lp.optionsRadioGroup.Changed() {
		lp.currentValue = lp.optionsRadioGroup.Value
		lp.WL.MultiWallet.SaveUserConfigValue(lp.preferenceKey, lp.optionsRadioGroup.Value)
		lp.updateButtonClicked()
		lp.RefreshTheme(lp.ParentWindow())
		lp.Dismiss()
	}

	for lp.cancelButton.Button.Clicked() {
		lp.Modal.Dismiss()
	}

	if lp.Modal.BackdropClicked(true) {
		lp.Modal.Dismiss()
	}
}

func (lp *ListPreferenceModal) Layout(gtx layout.Context) layout.Dimensions {
	w := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			txt := lp.Theme.H6(values.String(lp.title))
			txt.Color = lp.Theme.Color.Text
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.
				Layout(gtx, layout.Rigid(txt.Layout), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: values.MarginPaddingMinus2,
					}.Layout(gtx, lp.cancelButton.Layout)
				}))
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, lp.layoutItems()...)
		},
	}

	return lp.Modal.Layout(gtx, w)
}

func (lp *ListPreferenceModal) layoutItems() []layout.FlexChild {

	items := make([]layout.FlexChild, 0)
	for _, k := range lp.itemKeys {
		radioItem := layout.Rigid(lp.Theme.RadioButton(lp.optionsRadioGroup, k, values.String(lp.items[k]), lp.Theme.Color.DeepBlue, lp.Theme.Color.Primary).Layout)

		items = append(items, radioItem)
	}

	return items
}
