package constants

const (
	JAVA       = "java"
	PYTHON     = "python"
	GOLANG     = "golang"
	JAVASCRIPT = "javascript"
	CPP        = "c++"
)

const (
	ACCEPTED  = "Accepted"
	ATTEMPTED = "Attempted"
)

var IS_VALID_LANGUAGE = map[string]bool{
	PYTHON: true,
	//JAVA:       true,
	//GOLANG:     true,
	JAVASCRIPT: true,
	//CPP:        true,
}
