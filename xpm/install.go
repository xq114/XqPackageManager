package xpm

import (
	"fmt"
	"os"
	"path/filepath"
)

/*XPackage interface*/
type XPackage interface {
	GetRemoteVersion(string) (string, error)
	GetCurrentVersion(string) string
	UpdateFiles(string) error
}

/*Install p@v in configuration c at prefix*/
func Install(p string, v string, c []string, prefix string) error {
	switch p {
	case "imgui":
		xp := XImgui{config: nil}
		for _, cf := range c {
			xp.config = append(xp.config, cf)
		}
		fmt.Println("Try installing package Imgui...")
		prefix = filepath.Dir(prefix + "/imgui/")
		err := InstallPackage(&xp, v, prefix)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("Cannot find package %s", p)
	}
}

/*InstallPackage with version v*/
func InstallPackage(p XPackage, v string, prefix string) error {
	vsc := p.GetCurrentVersion(prefix)
	vs, err := p.GetRemoteVersion(v)
	if err != nil {
		return err
	}
	if vsc == vs {
		fmt.Println("The package is up to date.")
		return nil
	}
	fmt.Printf("Installing version %s\n", vs)
	os.RemoveAll(prefix)
	err = os.MkdirAll(prefix, 0755)
	if err != nil {
		return err
	}
	err = p.UpdateFiles(prefix)
	if err != nil {
		return err
	}
	return nil
}
