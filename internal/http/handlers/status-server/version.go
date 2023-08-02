package statusserver

import (
	"encoding/json"
	"net/http"

	"github.com/requiemofthesouls/logger"
)

type version struct {
	CommitHash string
	Branch     string
	Tag        string
	BuildDate  string
	BuiltBy    string
}

func GetVersionFromParams(params map[string]interface{}) *version {
	return &version{
		CommitHash: params["commit_hash"].(string),
		Branch:     params["branch"].(string),
		Tag:        params["tag"].(string),
		BuildDate:  params["build_date"].(string),
		BuiltBy:    params["built_by"].(string),
	}
}

func (m *manager) Version() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")

		var (
			content []byte
			err     error
		)
		if content, err = json.Marshal(m.version); err != nil {
			m.l.Error("error json.Marshal", logger.Error(err))
			http.Error(resp, "", http.StatusInternalServerError)
			return
		}

		if _, err = resp.Write(content); err != nil {
			m.l.Error("error resp.Write", logger.Error(err))
			http.Error(resp, "", http.StatusInternalServerError)
		}
	})
}
