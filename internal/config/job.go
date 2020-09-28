package config

// User can provide a local job config to instruct kaectl to prepare job
// It mainly used to upload artifacts
type JobConfig struct {
	Artifacts []string
}
