package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"rgx/common/log"
	"runtime"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func ShowConfig() {
	fmt.Println("Emvironment Variables:")
	printEnvVariable("RGX_CONFIG_DIR")
	cf := configFileName()
	fmt.Printf("\nConfig file: %v\n", cf)
}

func Configure() {
	var e error

	// set up directories we need
	e = os.MkdirAll(Config.DownloadDir, 0775)
	ErrCheck(e)
	e = os.MkdirAll(Config.PackagesDir, 0775)
	ErrCheck(e)
	e = os.MkdirAll(TempDir(), 0775)
	ErrCheck(e)
}

var ProgramSettings Dict

func ReadConfig() RgxConfig {
	configFile := configFileName()
	configBytes, e := os.ReadFile(configFile)
	if e != nil {
		fmt.Printf("error: couldn't open config file: '%v'\n", configFile)
		fmt.Println("    Ensure the file exists, and unset the RGX_CONFIG_DIR")
		fmt.Println("    environment variable if necessary.")
		os.Exit(1)
	}

	var config RgxConfig
	e = toml.Unmarshal(configBytes, &ProgramSettings)
	ErrCheck(e)

	plat := PlatformOS()
	config.ServerUrl = ProgramSettings.GetString("server_url", "")
	config.ShowProgress = ProgramSettings.GetBool("show_progress", false)

	config.PackagesDir = replaceTilde(ProgramSettings.GetDict(plat).GetString("packages_dir", "~/rgx-packages"))
	config.DownloadDir = replaceTilde(ProgramSettings.GetDict(plat).GetString("download_dir", "/tmp"))
	config.RcFileDir = replaceTilde(ProgramSettings.GetDict(plat).GetString("rcfile_dir", "~"))

	return config
}

func replaceTilde(s string) string {
	homeFolder := ""
	switch p := PlatformOS(); p {
	case "linux":
		homeFolder = os.Getenv("HOME")
	case "windows":
		homeFolder = os.Getenv("USERPROFILE")
	case "mac":
		homeFolder = os.Getenv("HOME")
	}
	if homeFolder != "" {
		return strings.Replace(s, "~", homeFolder, 1)
	} else {
		return s
	}
}

func printEnvVariable(v string) {
	r := os.Getenv(v)
	if r == "" {
		r = "[not set]"
	}
	fmt.Printf("  %v = %v\n", v, r)
}

func configFileName() string {
	var configFile = "rgx.toml"
	configDir := os.Getenv("RGX_CONFIG_DIR")
	if configDir != "" {
		fmt.Printf("reading config from: %s\n", configDir)
		configFile = filepath.Join(configDir, configFile)
	} else { // the config file is in the same directory as the exe
		exePath, e := os.Executable()
		ErrCheck(e)
		exeDir := filepath.Dir(exePath)
		configFile = filepath.Join(exeDir, configFile)
	}
	return configFile
}

//goland:noinspection GoBoolExpressions
func PlatformOS() string {
	p := runtime.GOOS
	if p == "darwin" {
		return "macos"
	} else {
		return p
	}
}

//goland:noinspection GoBoolExpressions
func PlatformArch() string {
	p := runtime.GOARCH
	if p == "amd64" {
		return "x64" // Devices with Intel 64-bit CPUs
	} else {
		return p
	}
}

func WinowsPathToUnix(path string) string {
	// only works on ully qualified paths, i.e. those that start with a drive letter, e.g. Z:
	return "/" + strings.ToLower(string(path[0])) + strings.ReplaceAll(path[2:], "\\", "/")
}

func Exists(fname string) bool {
	if _, e := os.Stat(fname); e == nil {
		return true
	} else {
		return false
	}
}

func ErrCheck(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

func TempDir() string {
	return filepath.Join(os.TempDir(), "rgx-temp")
}

func Assert(condition bool, failMsg string) {
	if !condition {
		panic(failMsg)
	}
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func CompareVersion(ver1, ver2 string) (int, error) {
	if ver1 == ver2 {
		return 0, nil
	}

	c1 := strings.Count(ver1, ".")
	c2 := strings.Count(ver2, ".")
	if c1 != c2 {
		return 0, errors.New("version strings are not comparable")
	}

	v1Parts := strings.Split(ver1, ".")
	v2Parts := strings.Split(ver2, ".")

	for i := 0; i < c1; i++ {
		v1, err := strconv.Atoi(v1Parts[i])
		if err != nil {
			return 0, err
		}
		v2, err := strconv.Atoi(v2Parts[i])
		if err != nil {
			return 0, err
		}
		if v1 < v2 {
			return -1, nil
		} else if v1 > v2 {
			return 1, nil
		}
	}

	return 0, errors.New("unexpected error")
}

func GetRuntimeConfig() RuntimeConfig {
	hostname, err := os.Hostname()
	if err != nil {
		log.Debug("could not get hostname: %s", err.Error())
		hostname = "motprovided"
	}
	currentUser, err := user.Current()
	var u string
	if err != nil {
		log.Debug("could not get current user: %s", err.Error())
		u = "notprovided"
	} else {
		u = currentUser.Username
	}
	opsys := PlatformOS()
	arch := PlatformArch()
	return RuntimeConfig{hostname, opsys, arch, u}
}

func (r *RuntimeConfig) AsHeader() string {
	return fmt.Sprintf("arch=%s os=%s machine=%s user=%s", r.Arch, r.OS, r.Machine, r.User)
}
