package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"rgx/common/log"
	"rgx/common/utils"
)

type TextResponse struct {
	Text         string
	ResponseCode int
}

func GetText(url string) (TextResponse, error) {
	client, req := setup(url, &utils.Config)
	resp, e := client.Do(req)
	if e != nil {
		return TextResponse{"", 0}, e
	}
	if resp.StatusCode != 200 {
		return TextResponse{"", resp.StatusCode}, errors.New(resp.Status)
	}
	respBody, e := io.ReadAll(resp.Body)
	if e != nil {
		return TextResponse{"", 0}, e
	}
	return TextResponse{string(respBody), 200}, nil
}

func Download(url, targetFile string, cksum utils.Checksum) (error, bool) {
	client, req := setup(url, &utils.Config)
	resp, e := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("could not close http response body: %s", err.Error())
		}
	}(resp.Body)
	utils.ErrCheck(e)

	var downloadSize uint64
	var showDownloadProgress bool
	if clen, ok := resp.Header["Content-Length"]; ok {
		downloadSize, e = strconv.ParseUint(clen[0], 10, 64)
		if e != nil {
			log.Warn("invalid value in content-length header")
		}
	}

	if downloadSize > 600000 {
		if utils.Config.ShowProgress {
			showDownloadProgress = true
		}
		log.Debug("about to download %v bytes\n", downloadSize)
	}

	tempFile := targetFile + ".rgxdownload"
	out, e := os.Create(tempFile)
	utils.ErrCheck(e)

	if showDownloadProgress {
		counter := &WriteCounter{TotalBytes: downloadSize}
		if _, e = io.Copy(out, io.TeeReader(resp.Body, counter)); e != nil {
			err := out.Close()
			if err != nil {
				log.Error("could not close downloaded file: %s", err.Error())
			}
		}
	}
	fmt.Printf("\r%s\r", strings.Repeat(" ", 40))

	_, e = io.Copy(out, resp.Body)
	if e != nil {
		return e, false
	}
	e = out.Close()
	if e != nil {
		return e, false
	}
	if utils.Exists(targetFile) {
		e = os.Remove(targetFile)
		if e != nil {
			return e, false
		}
	}
	e = os.Rename(tempFile, targetFile)
	if e != nil {
		return e, false
	}

	if cksum.Hash == "" {
		return nil, utils.Exists(targetFile)
	} else {
		s, e := utils.Hash(targetFile, cksum.Algorithm)
		if e != nil {
			return e, false
		}
		if s.Hash != strings.ToLower(strings.TrimSpace(cksum.Hash)) {
			log.Debug("checksum mismatch for: %s, expected %s but got %s", targetFile, cksum.Hash, s.Hash)
			return errors.New("checksum mismatch"), false
		} else {
			return nil, true
		}
	}
}

func SaveUrl(url, targetFile string) (error, bool) {
	client, req := setup(url, &utils.Config)
	resp, e := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("could not close http response body: %s", err.Error())
		}
	}(resp.Body)
	utils.ErrCheck(e)

	tempFile := targetFile + ".rgxdownload"
	out, e := os.Create(tempFile)
	utils.ErrCheck(e)

	_, e = io.Copy(out, resp.Body)
	if e != nil {
		return e, false
	}
	e = out.Close()
	if e != nil {
		return e, false
	}
	if utils.Exists(targetFile) {
		e = os.Remove(targetFile)
		if e != nil {
			return e, false
		}
	}
	e = os.Rename(tempFile, targetFile)
	if e != nil {
		return e, false
	}

	return nil, true
}

func setup(url string, config *utils.RgxConfig) (*http.Client, *http.Request) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent())
	if strings.HasPrefix(url, config.ArtifactRegistryBase) && config.ArtifactRegistryAuth != "" {
		log.Trace("adding well known credentials to request")
		u, p := credentials(config.ArtifactRegistryAuth)
		req.SetBasicAuth(u, p)
	}
	if strings.HasPrefix(url, config.ServerUrl) {
		req.Header.Set("x-rgx-installation", utils.CurrentRuntimeConfig.AsHeader())
	}
	return client, req
}

func userAgent() string {
	return utils.UserAgent
}

func credentials(auth string) (string, string) {
	creds := strings.Split(auth, ":")
	utils.Assert(len(creds) == 2, "unexpected number of credentials")
	return creds[0], creds[1]
}

type WriteCounter struct {
	BytesTransferred uint64
	TotalBytes       uint64
}

var writtenBytes = 0

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.BytesTransferred += uint64(n)
	writtenBytes += n
	if writtenBytes > 2000000 {
		wc.PrintProgress()
		writtenBytes = 0
	}
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	b := float64(wc.BytesTransferred)
	fmt.Printf("\r%s", strings.Repeat(" ", 40)) // clear line
	fmt.Printf("\rdownloading... %0.f MB complete ", b/1e6)
	if wc.TotalBytes > 0 {
		fmt.Printf("(%.0f%%)", b/float64(wc.TotalBytes)*100)
	}
}
