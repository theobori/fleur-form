package fleurform

import "strings"

const (
	ParametersBegin     = "\x01"
	ParametersSeparator = "\x02"
	PairSeparator       = "\x03"
)

type Parameters map[string]string

func GetParametersFromString(parametersString string) Parameters {
	parameters := Parameters{}

	parametersStringSplitted := strings.Split(parametersString, ParametersSeparator)
	for _, argumentString := range parametersStringSplitted {
		pair := strings.Split(argumentString, PairSeparator)
		if len(pair) != 2 {
			continue
		}

		parameters[pair[0]] = pair[1]
	}

	return parameters
}

func GetParametersFromPath(path string) Parameters {
	virtualPathSplitted := strings.Split(path, ParametersBegin)
	if len(virtualPathSplitted) < 2 {
		return Parameters{}
	}

	parametersString := virtualPathSplitted[1]

	return GetParametersFromString(parametersString)
}
