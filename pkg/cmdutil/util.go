package cmdutil

import (
	"fmt"
	"github.com/kaecloud/kaectl/api"
	"github.com/kaecloud/kaectl/pkg/spec"
	"github.com/kaecloud/kaectl/utils"
	"github.com/pkg/errors"
	"os"
	"path"
	"strings"
)

func PrepareJob(sp *spec.JobSpec, c *api.JobClient) (err error) {
	if sp.Prepare == nil || len(sp.Prepare.Artifacts) == 0 {
		return nil
	}

	for idx, artifact := range sp.Prepare.Artifacts {
		if artifact.Local == "" {
			continue
		}
		if path.IsAbs(artifact.Local) {
			return errors.Errorf("Only relative path is allowed, but got %s", artifact.Local)
		}
		var objUrl string
		zipFileName := fmt.Sprintf("%s.zip", strings.TrimRight(artifact.Local, " /\\"))
		err = utils.Compress(zipFileName, artifact.Local)
		if err != nil {
			return err
		}
		objUrl, err = c.UploadArtifact(sp.Name, zipFileName, zipFileName)
		os.Remove(zipFileName)

		if err != nil {
			return err
		}
		// Update artifact url
		sp.Prepare.Artifacts[idx].Url = objUrl
	}
	return nil
}
