package profile

type LoadProfileMode string

const (
	ModeConstantVUs      LoadProfileMode = "constant-vus"
	ModeSharedIterations LoadProfileMode = "shared-iterations"
	ModeStages           LoadProfileMode = "stages"
)

type StageConfig struct {
	Duration string
	Target   int
}

// Модель профиля нагрузки
type LoadProfile struct {
	Mode       LoadProfileMode
	InputPath  string
	OutputPath string
	VUs        int
	Duration   string
	Iterations int
	Stages     []StageConfig
	Thresholds map[string][]string
}

type BuildInput struct {
	InputPath     string
	OutputPath    string
	VUs           int
	Duration      string
	Iterations    int
	StagesRaw     []string
	ThresholdsRaw []string
}
