package projectview

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ProjectView struct {
	widget.BaseWidget
	Project  *binding.Struct
	ProjPath *binding.String
}

func NewProjectView(project *binding.Struct, projPath *binding.String) *ProjectView {
	w := &ProjectView{
		Project:  project,
		ProjPath: projPath,
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *ProjectView) CreateRenderer() fyne.WidgetRenderer {
	return newProjectViewRenderer(w)
}

type projectViewRenderer struct {
	widget     *ProjectView
	background *canvas.Rectangle
	name       *canvas.Text
	desc       *widget.Label
}

func newProjectViewRenderer(projectView *ProjectView) *projectViewRenderer {
	proj := projectView.Project
	name, e := (*proj).GetValue("Name")
	if e != nil {
		log.Fatal(e)
	}
	desc, e := (*proj).GetValue("Desc")
	if e != nil {
		log.Fatal(e)
	}
	title := cases.Title(language.Und)
	vr := &projectViewRenderer{
		widget:     projectView,
		background: canvas.NewRectangle(theme.BackgroundColor()),
		name:       canvas.NewText(title.String(name.(string)), theme.ForegroundColor()),
		desc:       widget.NewLabel(desc.(string)),
	}
	vr.name.TextStyle.Bold = true
	vr.name.TextSize = 32
	return vr
}

func (r *projectViewRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.name, r.desc}
}

func (r *projectViewRenderer) Layout(s fyne.Size) {
	ns := fyne.MeasureText(r.name.Text, r.name.TextSize, r.name.TextStyle)
	r.name.Move(fyne.Position{X: (s.Width - ns.Width) / 2, Y: 0})
	r.desc.Move(fyne.Position{X: (s.Width - r.desc.MinSize().Width) / 2, Y: ns.Height + 3})
	r.background.Resize(s)
}

func (r *projectViewRenderer) MinSize() fyne.Size {
	ts := fyne.MeasureText(r.name.Text, r.name.TextSize, r.name.TextStyle)
	return fyne.NewSize(ts.Width+theme.Padding()*5, ts.Height+theme.Padding()*5)
}

func (r *projectViewRenderer) Refresh() {
	proj := r.widget.Project
	name, e := (*proj).GetValue("Name")
	if e != nil {
		log.Fatal(e)
	}
	title := cases.Title(language.Und)
	r.name.Text = title.String(name.(string))
	r.name.Color = theme.ForegroundColor()
	r.background.FillColor = theme.BackgroundColor()
	r.background.Refresh() // Redraw the background first
	r.name.Refresh()       // Redraw the name on top
}

func (r *projectViewRenderer) Destroy() {}
