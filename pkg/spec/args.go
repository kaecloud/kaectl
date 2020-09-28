package spec

type CreateJobArgs struct {
	Name        string `json:"name,omitempty"`
	Shell       bool   `json:"shell,omitempty"`
	Image       string `json:"image,omitempty"`
	Command     string `json:"command,omitempty"`
	GPU         int    `json:"gpu,omitempty"`
	AutoRestart bool   `json:"auto_restart,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Spec        string `json:"spec,omitempty"`
	Cluster     string `json:"cluster"`
}

