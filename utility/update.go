package utility

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/ShrewdSpirit/credman/data"
)

type asset struct {
	BrowserDownloadUrl string `json:"browser_download_url"`
	ContentType        string `json:"content_type"`
	CreatedAt          string `json:"created_at"`
	DownloadCount      int    `json:"download_count"`
	Id                 int64  `json:"id"`
	Label              string `json:"label"`
	Name               string `json:"name"`
	NodeId             string `json:"node_id"`
	Size               int64  `json:"size"`
	State              string `json:"state"`
	UpdatedAt          string `json:"updated_at"`
	Uploader           author `json:"uploader"`
	Url                string `json:"url"`
}

type author struct {
	AvatarUrl         string `json:"avatar_url"`
	EventsUrl         string `json:"events_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	GravatarId        string `json:"gravatar_id"`
	HtmlUrl           string `json:"html_url"`
	Id                int64  `json:"id"`
	Login             string `json:"login"`
	NodeId            string `json:"node_id"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	ReposUrl          string `json:"repos_url"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	Type              string `json:"type"`
	Url               string `json:"url"`
}

type release struct {
	Assets          []asset `json:"assets"`
	AssetsUrl       string  `json:"assets_url"`
	Author          author  `json:"author"`
	Body            string  `json:"body"`
	CreatedAt       string  `json:"created_at"`
	Draft           bool    `json:"draft"`
	HtmlUrl         string  `json:"html_url"`
	Id              int64   `json:"id"`
	Name            string  `json:"name"`
	NodeId          string  `json:"node_id"`
	PreRelease      bool    `json:"prerelease"`
	PublishedAt     string  `json:"published_at"`
	TagName         string  `json:"tag_name"`
	TarballUrl      string  `json:"tarball_url"`
	TargetCommitish string  `json:"target_commitish"`
	UploadUrl       string  `json:"upload_url"`
	Url             string  `json:"url"`
	ZipballUrl      string  `json:"zipball_url"`
}

var (
	latestRelease = release{}
	assetIndex    = -1
)

func DueUpdateCheck() bool {
	if data.Config.AutoUpdateInterval == 0 {
		return false
	}
	lastUpdateCheck := time.Unix(data.Config.LastUpdateCheck, 0)
	return time.Now().Local().After(lastUpdateCheck.Add(time.Duration(data.Config.AutoUpdateInterval) * time.Hour * 24))
}

func CheckNewVersion() (string, error) {
	data.Config.LastUpdateCheck = time.Now().Local().Unix()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/ShrewdSpirit/credman/releases/latest", nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("request failed: %s", res.Status))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &latestRelease); err != nil {
		return "", err
	}

	for i, a := range latestRelease.Assets {
		extIndex := strings.Index(a.Name, ".")
		nameParts := strings.Split(a.Name[:extIndex], "-")[1:]
		if nameParts[0] == runtime.GOOS && nameParts[1] == runtime.GOARCH {
			assetIndex = i
			break
		}
	}

	if assetIndex == -1 {
		return "", errors.New(fmt.Sprintf("No release available for %s %s", runtime.GOOS, runtime.GOARCH))
	}

	return latestRelease.TagName, nil
}

func GetUpdate() error {
	if _, err := os.Stat(path.Join(data.DataDir, "origpath")); err == nil {
		return nil
	}

	if assetIndex == -1 {
		v, err := CheckNewVersion()
		if err != nil {
			return err
		}

		if v == data.Version {
			return errors.New("Up to date!")
		}
	}

	updateFilePath := path.Join(data.DataDir, "update.tar.gz")

	res, err := http.Get(latestRelease.Assets[assetIndex].BrowserDownloadUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	downloadFile, err := os.Create(updateFilePath)
	if err != nil {
		return err
	}
	defer downloadFile.Close()

	if _, err = io.Copy(downloadFile, res.Body); err != nil {
		return err
	}

	downloadFile.Seek(0, io.SeekStart)

	gzReader, err := gzip.NewReader(downloadFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF: // no more files to iterate
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		if header.Typeflag == tar.TypeReg {
			if strings.HasPrefix(header.Name, "credman") {
				outputFile, err := os.OpenFile(path.Join(data.DataDir, "update"), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					return err
				}
				defer outputFile.Close()

				if _, err = io.Copy(outputFile, tarReader); err != nil {
					return err
				}

				break
			}
		}
	}

	return nil
}

func InstallUpdate() error {
	exeFile, err := os.Executable()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(data.DataDir, "origpath"), []byte(exeFile), os.ModePerm); err != nil {
		return err
	}

	return ForkSelf(true, true, path.Join(data.DataDir, "update"), "update", "install")
}
