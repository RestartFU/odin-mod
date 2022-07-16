package dependency

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/RestartFU/gophig"
	"golang.org/x/exp/slices"
)

type Manager struct{}

func (m Manager) Init() error {
	path, err := exec.LookPath("odin.exe")
	if err != nil {
		return err
	}
	out, err := exec.Command(path, "version").CombinedOutput()
	if err != nil {
		return err
	}
	module := Module{
		Version:      "odin" + strings.Split(string(out), "version")[1],
		Dependencies: []string{},
	}
	return gophig.SetConfComplex("./odin.toml", gophig.TOMLMarshaler{}, module, 0664)
}
func (m Manager) DownloadDependencies(path string, force bool) error {
	mod, err := ReadModule(path)
	if err != nil {
		return err
	}
	odinPath, err := m.OdinPath()
	if err != nil {
		return err
	}
	for _, d := range mod.Dependencies {
		path := fmt.Sprintf("%s\\shared\\%s", odinPath, strings.ReplaceAll(d, "/", "_"))
		if _, err := os.Stat(path); !os.IsNotExist(err) && !force {
			_, err := ReadModule(path)
			if err != nil {
				continue
			}
			m.DownloadDependencies(path, false)
			continue
		}
		m.CloneRepository(d)
		_, err := ReadModule(path)
		if err != nil {
			continue
		}
		m.DownloadDependencies(path, force)
	}
	return nil
}

func (m Manager) CloneRepository(baseUrl string) error {
	mod, err := ReadModule(".")
	if err != nil {
		return err
	}
	var outputPath string
	odinPath, err := m.OdinPath()
	if err != nil {
		return err
	}
	branch := strings.Split(baseUrl, "@")
	var out []byte
	seed := rand.New(rand.NewSource(time.Now().Unix())).Int()
	tempFolder := fmt.Sprintf("%s\\temp%d", odinPath, seed)
	if len(branch) >= 2 {
		outputPath = strings.ToLower(strings.Split(fmt.Sprintf("%s\\shared\\%s", odinPath, strings.ReplaceAll(baseUrl, "/", "_")), "@")[0])
		url := "https://" + strings.Split(baseUrl, "@")[0]
		out, err = exec.Command("git.exe", "clone", "--branch", branch[1], url, "-output", tempFolder).CombinedOutput()
	} else {
		outputPath = strings.ToLower(fmt.Sprintf("%s\\shared\\%s", odinPath, strings.ReplaceAll(baseUrl, "/", "_")))
		url := "https://" + baseUrl
		out, err = exec.Command("git.exe", "clone", url, "-output", tempFolder).CombinedOutput()
	}
	defer m.DownloadDependencies(baseUrl, false)
	if err != nil {
		return errors.New(string(out))
	}
	os.RemoveAll(outputPath)
	os.Rename(tempFolder, outputPath)
	uri := baseUrl
	if len(branch) >= 2 {
		uri = strings.Split(baseUrl, "@")[0]
	}
	if !slices.Contains(mod.Dependencies, baseUrl) {
		mod.Dependencies = append(mod.Dependencies, uri)
	}
	return mod.Save()
}

func (Manager) OdinPath() (string, error) {
	path, err := exec.LookPath("odin.exe")
	if err != nil {
		return "", err
	}
	return strings.Split(path, "\\odin.exe")[0], nil
}
