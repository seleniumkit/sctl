package cmd

type Host struct {
	Port     int    `json:"port"`
	Count    int    `json:"count"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Region map[string]Host

type Regions map[string]Region

type Hosts map[string]Regions

type Versions map[string]string

type Browser struct {
	DefaultVersion string   `json:"defaultVersion"`
	Versions       Versions `json:"versions"`
}

type Quota map[string]Browser

type Input struct {
	Hosts   Hosts               `json:"hosts"`
	Quota   map[string]Quota    `json:"quota"`
	Aliases map[string][]string `json:"aliases"`
}
