package scanner

type BuildTool string

const (
	Maven  BuildTool = "MAVEN"
	Gradle BuildTool = "GRADLE"
)

type BuildWrapper struct {
	Windows string
	Unix    string
}

type Project struct {
	Name          string
	Path          string
	BuildTool     BuildTool
	BuildFile     string
	Wrapper       *BuildWrapper
	IsMultiModule bool
}
