package dtos

type GetPodDto struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	NodeName   string            `json:"nodename"`
	Phase      string            `json:"phase"`
	Labels     map[string]string `json:"labels"`
	Component  string            `json:"component"`
	Instance   string            `json:"instance"`
	Containers []string          `json:"containers"`
}
