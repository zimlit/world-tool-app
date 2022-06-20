package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	worldtoolprojects "github.com/zimlit/world-tool-projects"
	"github.com/zimlit/world-tool-projects/galaxy"
	projectview "github.com/zimlit/world-tool/projectView"
)

var proj binding.Struct
var projPath binding.String
var window fyne.Window

func main() {
	a := app.New()
	window = a.NewWindow("World Tool")

	// content

	// menu
	menuItem1 := fyne.NewMenuItem("New", newProj)
	menuItem2 := fyne.NewMenuItem("Save", saveProj)
	menuItem3 := fyne.NewMenuItem("Open", openProj)
	fileMenu := fyne.NewMenu("File", menuItem1, menuItem2, menuItem3)
	menu := fyne.NewMainMenu(fileMenu)
	window.SetMainMenu(menu)

	//shortcuts
	ctrlN := desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
	window.Canvas().AddShortcut(&ctrlN, func(shortcut fyne.Shortcut) { newProj() })
	ctrlS := desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	window.Canvas().AddShortcut(&ctrlS, func(shortcut fyne.Shortcut) { saveProj() })
	ctrlO := desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl}
	window.Canvas().AddShortcut(&ctrlO, func(shortcut fyne.Shortcut) { openProj() })

	window.ShowAndRun()
}

func saveProj() { fmt.Println("Save pressed") }

func newProj() {
	dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil {
			log.Fatal(err)
		}
		if uc == nil {
			return
		}
		uri := uc.URI()
		n := uri.Name()
		p := worldtoolprojects.New(n[:len(n)-len(filepath.Ext(n))], "", []galaxy.Galaxy{})
		projPath = binding.NewString()
		projPath.Set(uri.Path())
		entWidget := widget.NewEntry()
		entWidget.MultiLine = true
		ent := widget.NewFormItem("description", entWidget)
		dialog.ShowForm("Enter Description", "done", "cancel", []*widget.FormItem{ent}, func(b bool) {
			if b {
				p.Desc = entWidget.Text
				proj = binding.BindStruct(&p)
				v, e := xml.MarshalIndent(p, "", "  ")
				if e != nil {
					log.Fatal(e)
				}
				file, e := os.Create(uri.Path())
				if e != nil {
					log.Fatal(e)
				}
				_, e = file.Write(v)
				if e != nil {
					log.Fatal(e)
				}
				pv := projectview.NewProjectView(&proj, &projPath)
				window.SetContent(pv)
			}
		}, window)
	}, window)
}

func openProj() {
	dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err != nil {
			log.Fatal(err)
		}
		path := uc.URI().Path()

		projStr, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		var project worldtoolprojects.Project
		xml.Unmarshal(projStr, &project)
		if projPath == nil {
			projPath = binding.NewString()
		}
		projPath.Set(path)
		if proj == nil {
			proj = binding.BindStruct(&project)
		} else {
			proj.SetValue("Name", project.Name)
			proj.SetValue("Desc", project.Desc)
			proj.SetValue("Galaxies", project.Galaxies)
		}

		pv := projectview.NewProjectView(&proj, &projPath)
		window.SetContent(pv)
	}, window)
}
