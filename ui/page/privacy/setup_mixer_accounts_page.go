package privacy

import (
	"context"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/renderers"
	"github.com/planetdecred/godcr/ui/values"
)

const SetupMixerAccountsPageID = "SetupMixerAccounts"

type SetupMixerAccountsPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	wallet *dcrlibwallet.Wallet

	backButton              decredmaterial.IconButton
	infoButton              decredmaterial.IconButton
	autoSetupClickable      *decredmaterial.Clickable
	manualSetupClickable    *decredmaterial.Clickable
	autoSetupIcon, nextIcon *decredmaterial.Icon
}

func NewSetupMixerAccountsPage(l *load.Load, wallet *dcrlibwallet.Wallet) *SetupMixerAccountsPage {
	pg := &SetupMixerAccountsPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(SetupMixerAccountsPageID),
		wallet:           wallet,
	}
	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	pg.autoSetupIcon = decredmaterial.NewIcon(pg.Theme.Icons.ActionCheckCircle)
	pg.autoSetupIcon.Color = pg.Theme.Color.Success

	pg.nextIcon = decredmaterial.NewIcon(pg.Theme.Icons.NavigationArrowForward)
	pg.nextIcon.Color = pg.Theme.Color.Gray1

	pg.autoSetupClickable = pg.Theme.NewClickable(true)
	pg.manualSetupClickable = pg.Theme.NewClickable(true)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *SetupMixerAccountsPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *SetupMixerAccountsPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "Set up needed accounts",
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: func(gtx C) D {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									txt := pg.Theme.Body1("Two dedicated accounts will be set up to use the mixer:")
									txt.Alignment = text.Start
									ic := decredmaterial.NewIcon(pg.Theme.Icons.ImageBrightness1)
									ic.Color = pg.Theme.Color.Gray1
									return layout.Inset{Top: values.MarginPadding16, Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
										return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
											layout.Rigid(txt.Layout),
											layout.Rigid(func(gtx C) D {
												return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
													return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx C) D {
															return layout.Inset{Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
																return ic.Layout(gtx, values.MarginPadding8)
															})
														}),
														layout.Rigid(func(gtx C) D {
															txt2 := `<span style="text-color: grayText2">
														<b>Mixed </b> account will be the outbounding spending account.
													</span>`

															return layout.Inset{
																Left: values.MarginPadding8,
															}.Layout(gtx, renderers.RenderHTML(txt2, pg.Theme).Layout)
														}),
													)
												})
											}),
											layout.Rigid(func(gtx C) D {
												return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
													return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx C) D {
															return layout.Inset{Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
																return ic.Layout(gtx, values.MarginPadding8)
															})
														}),
														layout.Rigid(func(gtx C) D {
															txt3 := `<span style="text-color: grayText2">
													<b>Unmixed </b> account will be the change handling account.
												</span>`

															return layout.Inset{
																Left: values.MarginPadding8,
															}.Layout(gtx, renderers.RenderHTML(txt3, pg.Theme).Layout)
														}),
													)
												})
											}),
										)
									})
								}),
							)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return pg.autoSetupClickable.Layout(gtx, pg.autoSetupLayout)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return pg.manualSetupClickable.Layout(gtx, pg.manualSetupLayout)
						}),
					)
				})
			},
		}
		return page.Layout(pg.ParentWindow(), gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *SetupMixerAccountsPage) autoSetupLayout(gtx C) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.autoSetupIcon.Layout(gtx, values.MarginPadding20)
					}),
					layout.Rigid(func(gtx C) D {
						autoSetupText := pg.Theme.H6("Auto setup")
						txt := pg.Theme.Body2("Create and setup the needed accounts for you.")
						return layout.Inset{
							Left: values.MarginPadding16,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(autoSetupText.Layout),
								layout.Rigid(txt.Layout),
							)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Right: values.MarginPadding4,
					Top:   values.MarginPadding10,
				}.Layout(gtx, func(gtx C) D {
					return pg.nextIcon.Layout(gtx, values.MarginPadding20)
				})
			}),
		)
	})
}

func (pg *SetupMixerAccountsPage) manualSetupLayout(gtx C) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(pg.Theme.Icons.EditIcon.Layout24dp),
					layout.Rigid(func(gtx C) D {
						autoSetupText := pg.Theme.H6("Manual setup")
						txt := pg.Theme.Body2("For wallets that have enabled privacy before.")
						return layout.Inset{
							Left: values.MarginPadding16,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(autoSetupText.Layout),
								layout.Rigid(txt.Layout),
							)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Right: values.MarginPadding4,
					Top:   values.MarginPadding10,
				}.Layout(gtx, func(gtx C) D {
					return pg.nextIcon.Layout(gtx, values.MarginPadding20)
				})
			}),
		)
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SetupMixerAccountsPage) HandleUserInteractions() {
	if pg.autoSetupClickable.Clicked() {
		showModalSetupMixerInfo(&sharedModalConfig{
			Load:          pg.Load,
			window:        pg.ParentWindow(),
			pageNavigator: pg.ParentNavigator(),
			wallet:        pg.wallet,
			checkBox:      pg.Theme.CheckBox(new(widget.Bool), "Automatically move funds from default to unmixed account"),
		})
	}

	if pg.manualSetupClickable.Clicked() {
		pg.ParentNavigator().Display(NewManualMixerSetupPage(pg.Load, pg.wallet))
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *SetupMixerAccountsPage) OnNavigatedFrom() {
	pg.ctxCancel()
}
