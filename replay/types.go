package replay

type AuditEvent struct {
	Stage                    string `json:"stage"`
	RequestURI               string `json:"requestURI"`
	Verb                     string `json:"verb"`
	RequestReceivedTimestamp string `json:"requestReceivedTimestamp"`
	StageTimestamp           string `json:"stageTimestamp"`
	User                     struct {
		Username string   `json:"username"`
		Groups   []string `json:"groups"`
	} `json:"user"`
	ObjectRef struct {
		Resource   string `json:"resource"`
		Namespace  string `json:"namespace"`
		Name       string `json:"name"`
		APIVersion string `json:"apiVersion"`
	} `json:"objectRef"`
	SourceIPs      []string `json:"sourceIPs"`
	UserAgent      string   `json:"userAgent"`
	ResponseStatus struct {
		Code    int    `json:"code"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
	} `json:"responseStatus"`
}

// EventItem represents Kubernetes events
type EventItem struct {
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	InvolvedObject struct {
		Kind      string `json:"kind"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"involvedObject"`
}
