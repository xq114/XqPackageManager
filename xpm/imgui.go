package xpm

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"strconv"

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

	bigger_v180 := true
	subfolder := "backends"
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

		ver, err := strconv.ParseFloat(version, 64)
		if err == nil && ver < 1.80 {
			subfolder = "examples"
			bigger_v180 = false
		}
	}
	
	
	x.url = fmt.Sprintf("https://raw.githubusercontent.com/ocornut/imgui/%s", tag)
	x.filelist = append(x.filelist,
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
	)
	if bigger_v180 {
		x.filelist = append(x.filelist, "imgui_tables.cpp")
	}
	if len(x.config) == 0 {
		fmt.Println("using default config glfw & opengl3")
		x.config = append(x.config, "glfw", "opengl3")
	}
	for _, c := range x.config {
		switch c {
		case "vulkan":
			x.filelist = append(x.filelist, "vulkan/generate_spv.sh", "vulkan/glsl_shader.frag", "glsl_shader.vert")
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
			header := fmt.Sprintf("%s/imgui_impl_%s.h", subfolder, c)
			source := fmt.Sprintf("%s/imgui_impl_%s.cpp", subfolder, c)
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
