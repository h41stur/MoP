package messages

type NewAgent struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	SO       string `json:"so"`
	Arch     string `json:"arch"`
	Hostname string `json:"hostname"`
}

type PostCommandMessage struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
	Resp    string `json:"resp"`
}
