package candidates

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"rgx/common/http"
	"rgx/common/log"
	"rgx/common/utils"
)

func Install(pkg, suppliedMajorVersion string, lts bool) {
	var majorVersion string
	var e error
	var cleanDirs []string
	defer utils.CleanDirs(func() []string {
		return cleanDirs
	})

	if suppliedMajorVersion == "latest" {
		var versions = MajorVersions(pkg, lts)
		if len(versions) == 0 {
			log.Fatal("no versions found for package: %s", pkg)
		}
		majorVersion = versions[len(versions)-1]
	} else {
		majorVersion = suppliedMajorVersion
		utils.ErrCheck(e)
	}
	r := DownloadRecipe(pkg, majorVersion)

	for _, a := range r.Artifacts {

		target := filepath.Join(utils.Config.DownloadDir, a.Name)
		exists := utils.Exists(target)
		var sumOk = false
		if exists && a.Checksum != "" {
			sum, e := utils.Hash(target, a.ChecksumType)
			if e != nil {
				log.Fatal("could not verify checksum of %s", target)
			}
			sumOk = sum.Hash == a.Checksum
		} else {
			sumOk = true
		}
		if !sumOk {
			err := os.Remove(target)
			if err != nil {
				log.Fatal("could not remove previously downloaded file: %s", target)
			}
		}

		if !utils.Exists(target) {
			cs := utils.Checksum{Algorithm: a.ChecksumType, Hash: a.Checksum}
			log.Trace("downloading %s", a.Link)
			err, ok := http.Download(a.Link, target, cs)
			if !ok {
				if err != nil {
					log.Fatal("failed to download %s: %s", a.Link, err.Error())
				} else {
					log.Fatal("failed to download %s", a.Link)
				}
			}
		} else {
			log.Debug("already exists, not downloading again: %s", target)
		}

		switch a.Action {
		case "extract":
			extractdir := filepath.Join(utils.Config.PackagesDir, normalizedPath(a.ExtractDir))
			targetdir := filepath.Join(utils.Config.PackagesDir, normalizedPath(a.ExtractTarget))
			if !utils.Exists(targetdir) {
				e = os.MkdirAll(extractdir, 0775)
				log.Info("extracting files to %s", extractdir)
				utils.Extract(target, extractdir)
				log.Debug("extracted to %s", extractdir)
				if !utils.Exists(targetdir) {
					log.Info("finished extracting to %s", targetdir)
				}
			} else {
				log.Info("%s already exists", targetdir)
			}
		case "extract-to-temp":
			dir, err := utils.ExtractToTemp(target, a.ArtifactType)
			if err != nil {
				log.Fatal("failed to extract %s: %s", target, err.Error())
			}
			cleanDirs = append(cleanDirs, dir)
			log.Debug("extracted to %s", dir)
		case "copy":
			targetfile := filepath.Join(utils.Config.PackagesDir, normalizedPath(a.ExtractTarget), a.Name)
			_, copyErr := utils.Copy(target, targetfile)
			if copyErr != nil {
				log.Fatal("failed to write %s: %s", a.Name, copyErr.Error())
			}
		}
	}

	if r.Script == "" || r.ScriptDir == "" {
		return
	}

	scriptUrl := utils.Config.ServerUrl + r.Script
	scriptBase := filepath.Base(r.Script)
	scriptDir := filepath.Join(utils.Config.PackagesDir, normalizedPath(r.ScriptDir))
	packageVersion := r.PackageVersion
	err, _ := http.SaveUrl(scriptUrl, filepath.Join(scriptDir, scriptBase))
	if err != nil {
		log.Fatal("failed to download %s", scriptUrl)
	}

	var envmap map[string]string
	envmap = make(map[string]string)
	envmap["RGX_PACKAGE_MAJORVERSION"] = majorVersion
	envmap["RGX_PACKAGE_SCRIPTDIR"] = normalizedPath(scriptDir)
	envmap["RGX_PACKAGE_SCRIPTDIR_MSYS"] = msysPath(utils.Config.PackagesDir)
	envmap["RGX_PACKAGES_DIR"] = utils.Config.PackagesDir
	envmap["RGX_PACKAGES_DIR_MSYS"] = msysPath(utils.Config.PackagesDir)
	envmap["RGX_RCFILE_DIR"] = utils.Config.RcFileDir
	envmap["RGX_PACKAGE_VERSION"] = packageVersion
	utils.RunScript(scriptBase, scriptDir, envmap)
}

func normalizedPath(p string) string {
	if utils.PlatformOS() == "windows" {
		return strings.ReplaceAll(p, "/", "\\")
	} else {
		return p
	}
}

func msysPath(p string) string {
	if utils.PlatformOS() == "windows" {
		t := strings.ReplaceAll(p, "\\", "/")
		t = strings.ReplaceAll(t, ":", "")
		return "/" + t
	} else {
		return p
	}
}

type recipe struct {
	Script         string `json:"script"`
	ScriptDir      string `json:"script_dir"`
	PackageVersion string `json:"package_version"`
	Artifacts      []struct {
		ArtifactType  string `json:"artifact_type"`
		Action        string `json:"action"`
		Name          string `json:"name"`
		Version       string `json:"version"`
		Link          string `json:"link"`
		Checksum      string `json:"checksum"`
		ChecksumType  string `json:"checksum_type"`
		ExtractDir    string `json:"extract_dir"`
		ExtractTarget string `json:"extract_target"`
	} `json:"artifacts"`
}

func DownloadRecipe(pkg, majorVersion string) recipe {
	var u = "/packages/" + pkg + "/release/" + majorVersion + "/" + utils.PlatformOS() + "/" + utils.PlatformArch()
	log.Debug("getting package details from %s", utils.Config.ServerUrl+u)
	resp, e := http.GetText(utils.Config.ServerUrl + u)
	if e != nil {
		if resp.ResponseCode == 404 {
			log.Fatal("package not found %s version %s (%s%s)", pkg, majorVersion, utils.PlatformOS(), utils.PlatformArch())
		} else {
			log.Fatal("invalid server response: %s", e.Error())
		}
	}

	var r recipe
	err := json.Unmarshal([]byte(resp.Text), &r)
	if err != nil {
		log.Fatal("could not parse json: " + err.Error())
	}
	return r
}
