// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/values"
)

type Switch struct {
	style    *values.SwitchStyle
	disabled bool
	value    bool
	changed  bool
	clk      *widget.Bool
}

type SwitchItem struct {
	Text   string
	button Button
}

type SwitchButtonText struct {
	t                                  *Theme
	ActiveTextColor, InactiveTextColor color.NRGBA
	Active, Inactive                   color.NRGBA
	items                              []SwitchItem
	selected                           int
	changed                            bool
}

func (t *Theme) Switch() *Switch {
	return &Switch{
		clk:   new(widget.Bool),
		style: t.Styles.SwitchStyle,
	}
}

func (t *Theme) SwitchButtonText(i []SwitchItem) *SwitchButtonText {
	sw := &SwitchButtonText{
		t:     t,
		items: make([]SwitchItem, len(i)+1),
	}
	sw.Active, sw.Inactive = sw.t.Color.Surface, color.NRGBA{}
	sw.ActiveTextColor, sw.InactiveTextColor = sw.t.Color.GrayText1, sw.t.Color.Text

	for index := range i {
		i[index].button = t.Button(i[index].Text)
		i[index].button.HighlightColor = t.Color.SurfaceHighlight
		i[index].button.Background, i[index].button.Color = sw.Inactive, sw.InactiveTextColor
		i[index].button.TextSize = unit.Sp(14)
		sw.items[index+1] = i[index]
	}

	if len(sw.items) > 0 {
		sw.selected = 1
	}
	return sw
}

func (s *Switch) Layout(gtx layout.Context) layout.Dimensions {
	dGtx := gtx

	trackWidth := dGtx.Dp(32)
	trackHeight := dGtx.Dp(20)
	thumbSize := dGtx.Dp(18)
	trackOff := (thumbSize - trackHeight) / 2

	// Draw track.
	trackCorner := trackHeight / 2
	trackRect := image.Rectangle{Max: image.Point{
		X: trackWidth,
		Y: trackHeight,
	}}

	activeColor, InactiveColor, thumbColor := s.style.ActiveColor, s.style.InactiveColor, s.style.ThumbColor
	if s.disabled {
		dGtx = gtx.Disabled()
		activeColor, InactiveColor, thumbColor = Disabled(activeColor), Disabled(InactiveColor), Disabled(thumbColor)
	}

	col := InactiveColor
	if s.IsChecked() {
		col = activeColor
	}

	trackColor := col
	t := op.Offset(image.Point{Y: trackOff}).Push(dGtx.Ops)
	cl := clip.UniformRRect(trackRect, trackCorner).Push(dGtx.Ops)
	paint.ColorOp{Color: trackColor}.Add(dGtx.Ops)
	paint.PaintOp{}.Add(dGtx.Ops)
	cl.Pop()
	t.Pop()

	// Compute thumb offset and color.
	if s.IsChecked() {
		xoff := trackWidth - thumbSize
		defer op.Offset(image.Point{X: xoff}).Push(dGtx.Ops).Pop()
	}

	thumbRadius := thumbSize / 2

	circle := func(x, y, r int) clip.Op {
		b := image.Rectangle{
			Min: image.Pt(x-r, y-r),
			Max: image.Pt(x+r, y+r),
		}
		return clip.Ellipse(b).Op(dGtx.Ops)
	}

	// Draw thumb shadow, a translucent disc slightly larger than the
	// thumb itself.
	// Center shadow horizontally and slightly adjust its Y.
	paint.FillShape(dGtx.Ops, col, circle(thumbRadius, thumbRadius+dGtx.Dp(.25), thumbRadius+1))

	// Draw thumb.
	paint.FillShape(dGtx.Ops, thumbColor, circle(thumbRadius, thumbRadius, thumbRadius))

	// Set up click area.
	clickSize := dGtx.Dp(38)
	clickOff := image.Point{
		X: (trackWidth) - (clickSize),
		Y: (trackHeight) - (clickSize)/2 + trackOff,
	}
	defer op.Offset(clickOff).Push(dGtx.Ops).Pop()
	sz := image.Pt(clickSize, clickSize)
	defer clip.Ellipse(image.Rectangle{Max: sz}).Push(dGtx.Ops).Pop()
	s.clk.Layout(dGtx, func(dGtx layout.Context) layout.Dimensions {
		semantic.Switch.Add(dGtx.Ops)
		return layout.Dimensions{Size: sz}
	})

	dims := image.Point{X: trackWidth, Y: thumbSize}
	return layout.Dimensions{Size: dims}
}

func (s *Switch) Changed() bool {
	return s.clk.Changed()
}

func (s *Switch) IsChecked() bool {
	return s.clk.Value
}

func (s *Switch) SetChecked(value bool) {
	s.clk.Value = value
}

func (s *Switch) SetEnabled(value bool) {
	s.disabled = value
}

func (s *SwitchButtonText) Layout(gtx layout.Context) layout.Dimensions {
	s.handleClickEvent()
	m8 := unit.Dp(8)
	m4 := unit.Dp(4)
	card := s.t.Card()
	card.Color = s.t.Color.Gray2
	card.Radius = Radius(8)
	return card.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
			list := &layout.List{Axis: layout.Horizontal}
			Items := s.items[1:]
			return list.Layout(gtx, len(Items), func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx C) D {
					index := i + 1
					btn := s.items[index].button
					btn.Inset = layout.Inset{
						Left:   m8,
						Bottom: m4,
						Right:  m8,
						Top:    m4,
					}
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(btn.Layout),
					)
				})
			})
		})
	})
}

func (s *SwitchButtonText) handleClickEvent() {
	for index := range s.items {
		if index != 0 {
			if s.items[index].button.Clicked() {
				if s.selected != index {
					s.changed = true
				}
				s.selected = index
			}
		}

		if s.selected == index {
			s.items[s.selected].button.Background = s.Active
			s.items[s.selected].button.Color = s.ActiveTextColor
		} else {
			s.items[index].button.Background = s.Inactive
			s.items[index].button.Color = s.InactiveTextColor
		}
	}
}

func (s *SwitchButtonText) SelectedOption() string {
	return s.items[s.selected].Text
}

func (s *SwitchButtonText) SelectedIndex() int {
	return s.selected
}

func (s *SwitchButtonText) Changed() bool {
	changed := s.changed
	s.changed = false
	return changed
}
