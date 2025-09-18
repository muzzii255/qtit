package dashboard

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Torrent struct {
	Name         string  `json:"name"`
	Progress     float64 `json:"progress"`
	State        string  `json:"state"`
	Speed        int     `json:"dlspeed"`
	UpSpeed      int     `json:"upspeed"`
	ETA          int     `json:"eta"`
	Peers        int     `json:"peers"`
	Size         int     `json:"size"`
	AddedOn      int     `json:"added_on"`
	Hash         string  `json:"hash"`
	Seeds        int     `json:"num_seeds"`
	Leech        int     `json:"num_leechs"`
	Private      bool    `json:"private"`
	ForceStart   bool    `json:"force_start"`
	SuperSeeding bool    `json:"super_seeding"`
}

type TorrentFile struct {
	Index    int     `json:"index"`
	Name     string  `json:"name"`
	Size     int     `json:"size"`
	Progress float64 `json:"progress"`
	Priority int     `json:"priority"`
	IsSeed   *bool   `json:"is_seed,omitempty"`
}

type Qbit struct {
	Url      string
	Username string
	Password string
}

func Logout(client *http.Client, host string) {
	url := fmt.Sprintf("%s/api/v2/auth/logout", strings.TrimRight(host, "/"))
	_, err := client.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func LoginToQbit(host, username, password string) (*http.Client, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	data := url.Values{
		"username": {username},
		"password": {password},
	}
	resp, err := client.PostForm(host+"/api/v2/auth/login", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return client, nil
}

func AddMagnet(client *http.Client, host, magnet string) error {
	addURL := fmt.Sprintf("%s/api/v2/torrents/add", strings.TrimRight(host, "/"))
	form := url.Values{}
	form.Set("urls", magnet)
	values := url.Values{"urls": {magnet}}
	resp, err := client.PostForm(addURL, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func FetchTorrents(client *http.Client, host string) ([]Torrent, error) {
	url := fmt.Sprintf("%s/api/v2/torrents/info?sort=ratio", strings.TrimRight(host, "/"))
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var torrents []Torrent
	err = json.Unmarshal(body, &torrents)

	if resp.StatusCode != 200 {
		return []Torrent{}, fmt.Errorf("Status %d: %s", resp.StatusCode, string(body))
	}
	return torrents, nil
}

func FetchTorrentFiles(client *http.Client, host, hash string) ([]TorrentFile, error) {
	url := fmt.Sprintf("%s/api/v2/torrents/files?hash=%s", strings.TrimRight(host, "/"), hash)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return []TorrentFile{}, fmt.Errorf("Status %d: %s", resp.StatusCode, string(body))
	}
	var files []TorrentFile
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func PostTorrentAction(client *http.Client, host, endpoint string, params url.Values) error {
	url := fmt.Sprintf("%s/api/v2/torrents/%s", strings.TrimRight(host, "/"), endpoint)
	resp, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}
	return nil
}
