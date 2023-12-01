package config

import "time"

// CodeList code list
type CodeList struct {
	Types         string `json:"types"`
	Title         string `json:"title"`
	Address       string `json:"address"`
	Port          int    `json:"port"`
	Password      string `json:"password"`
	Encryption    string `json:"encryption"`
	Security      string `json:"security"`
	HeaderType    string `json:"headerType"`
	Method        string `json:"method"`
	Protocol      string `json:"protocol"`
	ProtocolParam string `json:"protocol_param"`
	Obfs          string `json:"obfs"`
	ObfsParam     string `json:"obfs_param"`
	Net           string `json:"net"`
	Host          string `json:"host"`
	Path          string `json:"path"`
	TLS           bool   `json:"tls"`
	Aid           int    `json:"aid"`
	Ping          string `json:"ping"`
}

// SS ss
type SS struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Method   string `json:"method"`
	Password string `json:"password"`
	Types    string `json:"types"`
	Title    string `json:"title"`
}

// SSR ssr
type SSR struct {
	Server        string `json:"server"`
	ServerPort    int    `json:"server_port"`
	LocalAddress  string `json:"local_address"`
	LocalPort     int    `json:"local_port"`
	Timeout       int    `json:"timeout"`
	Password      string `json:"password"`
	Method        string `json:"method"`
	Protocol      string `json:"protocol"`
	ProtocolParam string `json:"protocol_param"`
	Obfs          string `json:"obfs"`
	ObfsParam     string `json:"obfs_param"`
}

// Vary is v2ray json sturct
type Vary struct {
	Host  string `json:"host"`
	Path  string `json:"path"`
	TLS   bool   `json:"tls"`
	Ps    string `json:"ps"`
	Add   string `json:"add"`
	Port  int    `json:"port"`
	ID    string `json:"id"`
	Aid   int    `json:"aid"`
	Net   string `json:"net"`
	Type  string `josn:"type"`
	Types string `json:"types"`
	Title string `json:"title"`
}

// Vless vless
type Vless struct {
	UUID       string `json:"uuid"`
	URL        string `json:"url"`
	Port       int    `json:"port"`
	Encryption string `json:"encryption"`
	Security   string `json:"security"`
	Type       string `json:"type"`
	HeaderType string `json:"headerType"`
	Types      string `json:"types"`
	Title      string `json:"title"`
}

// Node node
type Node struct {
	NODE string `form:"node" json:"node" xml:"node"  binding:"required"`
}

// Domains
type Domains struct {
	Proxy  string `form:"proxy" json:"proxy" xml:"proxy"  binding:"required"`
	Direct string `form:"direct" json:"direct" xml:"direct"  binding:"required"`
}

// Subscribes
type Subscribes struct {
	Subscribes string `form:"subscribes" json:"subscribes" xml:"subscribes"  binding:"required"`
}

// Ignore
type Ignore struct {
	Ignore string `form:"ignore" json:"ignore" xml:"ignore"`
}

// Dns
type Dns struct {
	Dns string `form:"dns" json:"dns" xml:"dns"`
}

// Socks
type Socks struct {
	SockStuts string `form:"sockStuts" json:"sockStuts" xml:"sockStuts"`
	Auths     string `form:"auths" json:"auths" xml:"auths"`
}

// Uri
type Uri struct {
	Uri string `form:"uri" json:"uri" xml:"uri" binding:"required"`
}

// Config config
type Config struct {
	DefaultSubUrl     string   `json:"default_sub_url"`
	GetCoreVersionUrl string   `json:"get_core_version_url"`
	CoreVersion       string   `json:"core_version"`
	ProxyUrl          []string `json:"proxy_url"`
	Current           string   `json:"current"`
	GeoVersionUrl     string   `json:"geo_version_url"`
	GeoVersion        string   `json:"geo_version"`
}

// Template template
type Template struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Config string `json:"config"`
}

// Trojan template
type Trojan struct {
	RunType    string   `json:"run_type"`
	LocalAddr  string   `json:"local_addr"`
	LocalPort  int      `json:"local_port"`
	RemoteAddr string   `json:"remote_addr"`
	RemotePort int      `json:"remote_port"`
	Password   []string `json:"password"`
	Ssl        Ssl      `json:"ssl"`
	Mux        Mux      `json:"mux"`
	Router     Router   `json:"router"`
}

// Ssl to trojan
type Ssl struct {
	Sni string `json:"sni"`
}

// Mux to trojan
type Mux struct {
	Enabled bool `json:"enabled"`
}

// Router to trojan
type Router struct {
	Enabled       bool     `json:"enabled"`
	Bypass        []string `json:"bypass"`
	Block         []string `json:"block"`
	Proxy         []string `json:"proxy"`
	DefaultPolicy string   `json:"default_policy"`
	Geoip         string   `json:"geoip"`
	Geosite       string   `json:"geosite"`
}

// Config
type Configs struct {
	DefaultSubUrl     string   `json:"default_sub_url"`
	GetCoreVersionUrl string   `json:"get_core_version_url"`
	ProxyUrl          []string `json:"proxy_url"`
	Current           string   `json:"current"`
}

type GopherInfo struct {
	ID string
}

type Message struct {
	Type string `json:"type"`
	UUID string `json:"uuid"`
	Data string `json:"data"`
}

type JSONData struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
	Reactions  struct {
		URL        string `json:"url"`
		TotalCount int    `json:"total_count"`
		Num1       int    `json:"+1"`
		Num10      int    `json:"-1"`
		Laugh      int    `json:"laugh"`
		Hooray     int    `json:"hooray"`
		Confused   int    `json:"confused"`
		Heart      int    `json:"heart"`
		Rocket     int    `json:"rocket"`
		Eyes       int    `json:"eyes"`
	} `json:"reactions"`
	MentionsCount int `json:"mentions_count"`
}
