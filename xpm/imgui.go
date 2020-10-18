package xpm

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/antchfx/xmlquery"
)

// XImgui package script
type XImgui struct {
	url      string
	filelist []string
	config   []string
}

// GetRemoteVersion of XImgui
func (x *XImgui) GetRemoteVersion(version string) (string, error) {
	var tag string = "master"
	tagsurl := "https://github.com/ocornut/imgui/refs-tags/master?source_action=disambiguate&source_controller=files&tag_name=master&q="
	res, err := GetResponseWithHeader(tagsurl)
	if err != nil {
		return "", err
	}

	doc, err := xmlquery.Parse(strings.NewReader(res))
	if err != nil {
		return "", err
	}
	if version == "latest" {
		node, err := xmlquery.Query(doc, "//span[contains(@class, \"css-truncate\")]/text()")
		if err != nil {
			return "", err
		}
		if node == nil {
			return "", fmt.Errorf("Cannot parse for latest version")
		}
		tag = strings.TrimSpace(node.Data)
	} else if version != "master" {
		node, err := xmlquery.Query(doc, fmt.Sprintf("//span[@title=\"v%s\"]/text()", version))
		if err != nil {
			return "", err
		}
		if node == nil {
			return "", fmt.Errorf("Cannot find version %s", version)
		}
		tag = strings.TrimSpace(node.Data)
	}

	x.url = fmt.Sprintf("https://raw.githubusercontent.com/ocornut/imgui/%s", tag)
	x.filelist = []string{
		"imconfig.h",
		"imgui.h",
		"imgui.cpp",
		"imgui_internal.h",
		"imgui_demo.cpp",
		"imgui_draw.cpp",
		"imgui_widgets.cpp",
		"imstb_rectpack.h",
		"imstb_textedit.h",
		"imstb_truetype.h",
	}
	for _, c := range x.config {
		switch c {
		case "glut":
			fmt.Fprintf(os.Stderr, "Warning: glut is obsolete, try using glfw or sdl instead!\n")
			fallthrough
		case "glfw":
			fallthrough
		case "sdl":
			fallthrough
		case "win32":
			fallthrough
		case "opengl2":
			fallthrough
		case "opengl3":
			fallthrough
		case "dx11":
			fallthrough
		case "dx12":
			fallthrough
		case "vulkan":
			header := fmt.Sprintf("examples/imgui_impl_%s.h", c)
			source := fmt.Sprintf("examples/imgui_impl_%s.cpp", c)
			x.filelist = append(x.filelist, header, source)
		default:
			fmt.Fprintf(os.Stderr, "Warning: config %s is not recognized\n", c)
		}
	}

	return tag, nil
}

// GetCurrentVersion of XImgui
func (x *XImgui) GetCurrentVersion(prefix string) string {
	if len(x.config) != 0 {
		return ""
	}
	files, err := ioutil.ReadDir(prefix)
	if err != nil || len(files) < 14 {
		return ""
	}
	filename := filepath.Join(prefix, "imgui.h")
	fp, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	scanner.Scan()
	vers := scanner.Text()[15:]
	if len(x.config) == 0 {
		fmt.Println("using default config glfw & opengl3")
		x.config = append(x.config, "glfw", "opengl3")
	}
	return vers
}

// UpdateFiles of XImgui
func (x XImgui) UpdateFiles(prefix string) error {
	for _, pfile := range x.filelist {
		nurl := fmt.Sprintf("%s/%s", x.url, pfile)
		nfile := filepath.Base(pfile)
		err := DownloadAndSave(nurl, prefix, nfile)
		if err != nil {
			return err
		}
	}
	return nil
}
