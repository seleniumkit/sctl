package cmd

type Host struct {
	Port     int    `json:"port"`
	Ports    string `json:"ports"`
	Count    int    `json:"count"`
	Username string `json:"username"`
	Password string `json:"password"`
	VNC      string `json:"vnc"`
}

type Region map[string]Host

type Regions map[string]Region

type Hosts map[string]Regions

type Versions map[string]string

type Browser struct {
	DefaultVersion  string   `json:"defaultVersion"`
	DefaultPlatform string   `json:"defaultPlatform"`
	Versions        Versions `json:"versions"`
}

type Quota map[string]Browser

type Input struct {
	Hosts   Hosts               `json:"hosts"`
	Quota   map[string]Quota    `json:"quota"`
	Aliases map[string][]string `json:"aliases"`
}
