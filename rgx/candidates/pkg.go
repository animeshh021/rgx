package candidates

import (
	"encoding/json"
	"fmt"
	"rgx/common/http"
	"rgx/common/log"
	"rgx/common/utils"
)

type packagesResponse struct {
	Name        string
	Description string
}

func PrintServerPackages() {
	resp, e := http.GetText(utils.Config.ServerUrl + "/packages")
	if e != nil {
		log.Fatal("could not connect to rgx server: %s", e.Error())
	}
	var pkgs []packagesResponse
	err := json.Unmarshal([]byte(resp.Text), &pkgs)
	if err != nil {
		log.Fatal("could not parse server json: " + err.Error())
	}

	for _, r := range pkgs {
		fmt.Printf("%s - %s\n", r.Name, r.Description)
	}
}

func PrintMajorVersions(pkg string, ltsOnly bool) {
	var versions = MajorVersions(pkg, ltsOnly)
	for _, r := range versions {
		fmt.Printf("%s ", r)
	}
	fmt.Println()
}

func MajorVersions(pkg string, ltsOnly bool) []string {
	var u = "/packages/" + pkg + "/versions"
	if ltsOnly {
		u += "?lts=1"
	}
	resp, e := http.GetText(utils.Config.ServerUrl + u)
	if e != nil {
		if resp.ResponseCode == 404 {
			log.Fatal("package not found: %s", pkg)
		} else {
			log.Fatal("invalid server response: " + e.Error())
		}
	}

	var releases []string
	err := json.Unmarshal([]byte(resp.Text), &releases)
	if err != nil {
		log.Fatal("could not parse server json: " + err.Error())
	}
	return releases
}
