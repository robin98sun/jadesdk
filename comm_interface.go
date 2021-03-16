package jadesdk

type CommScheme string

const (
	CommSchemeHTTP CommScheme = "http"
	CommSchemeTCP             = "tcp"
	CommSchemeUDP             = "udp"
)

type Comm struct {
	Scheme       CommScheme `json:"scheme,ommitempty"`
	RetryEnabled bool       `json:"retryEnabled,omitempty"`
	RetryLimit   int        `json:"retryLimit,omitempty"`
	Timeout      int        `json:"timeout,omitempty"` // in microseconds
	StateEnabled bool       `json:"stateEnabled,omitempty"`
}

type CommParams struct {
	Protocol string `json:"protocol,omitempty"`
	Method   string `json:"method,omitempty"`
	Path     string `json:"path,omitempty"`
}

func (c *Comm) Request(
	operationName string,
	targetNode *Node,
	params *CommParams,
	payload interface{}) (interface{}, int, error) {

	if c.Scheme == CommSchemeHTTP {

	}

	return nil, 0, nil
}
