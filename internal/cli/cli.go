package cli

import (
	"fmt"
	"strings"

	"github.com/elecbug/regdir/internal/cli/parameter"
	"github.com/elecbug/regdir/internal/manager"
)

// Run is the main entry point for the CLI application. It takes a Parameters struct and returns a
// result string and a boolean indicating success.
func Run(params *parameter.Parameters) (string, bool) {
	switch strings.ToLower(params.MainParam) {
	case HELP_CMD:
		return helpStr(params), true
	case VERSION_CMD:
		return versionStr(params), true
	case RENAME_CMD:
		root, mode, src, dst, overwrite, err := getRenameParameters(params)
		if err != nil {
			return fmt.Sprintf("Failed to get rename parameters: %v", err), false
		}

		switch strings.ToLower(mode) {
		case WILDCARD_MODE, WILDCARD_MODE_SHORT:
			err := manager.ChangeFileNamesWithPattern(root, src, dst, manager.WILDCARD, overwrite)
			if err != nil {
				return fmt.Sprintf("Failed to rename files: %v", err), false
			}

			return fmt.Sprintf("Files renamed successfully using wildcard pattern: %s -> %s", src, dst), true
		case REGEX_MODE, REGEX_MODE_SHORT:
			err := manager.ChangeFileNamesWithPattern(root, src, dst, manager.REGEX, overwrite)
			if err != nil {
				return fmt.Sprintf("Failed to rename files: %v", err), false
			}

			return fmt.Sprintf("Files renamed successfully using regex pattern: %s -> %s", src, dst), true
		default:
			return fmt.Sprintf("Unknown mode: %s", mode), false
		}
	default:
		fmt.Printf("Unknown command: %s\n", params.MainParam)
		return helpStr(params), false
	}
}

// getRenameParameters retrieves the parameters needed for the rename command from the Parameters struct.
func getRenameParameters(params *parameter.Parameters) (root, mode, src, dst string, overwrite bool, err error) {
	root, err = getParameter(params, ROOT_PARAM, ROOT_PARAM_SHORT)
	if err != nil {
		return
	}
	mode, err = getParameter(params, MODE_PARAM, MODE_PARAM_SHORT)
	if err != nil {
		return
	}
	src, err = getParameter(params, SRC_PARAM, SRC_PARAM_SHORT)
	if err != nil {
		return
	}
	dst, err = getParameter(params, DST_PARAM, DST_PARAM_SHORT)
	if err != nil {
		return
	}
	overwrite = params.Params[OVERWRITE_PARAM] == parameter.CLI_TRUE_STR

	if root == "" || mode == "" || src == "" || dst == "" {
		err = fmt.Errorf("missing required parameters for rename command")
		return
	}

	return
}

// getParameter is a helper function to retrieve a parameter value from the Parameters struct using multiple possible keys.
// It returns the value of the first key found and an error if none of the keys are found or if multiple keys are found.
func getParameter(params *parameter.Parameters, keys ...string) (string, error) {
	onlyOne := false
	result := ""
	for _, key := range keys {
		if value, exists := params.Params[key]; exists {
			if onlyOne {
				return "", fmt.Errorf("multiple keys found in parameters: %v", keys)
			}
			result = value
			onlyOne = true
		}
	}
	if !onlyOne {
		return "", fmt.Errorf("none of the keys %v found in parameters", keys)
	}
	return result, nil
}

// helpStr returns a usage string for the CLI application.
func helpStr(params *parameter.Parameters) string {
	result := `Usage: ` + params.MainParam + ` [OPTIONS]...`
	return result
}

// versionStr returns a version string for the CLI application, including the version number and author.
func versionStr(params *parameter.Parameters) string {
	result := fmt.Sprintf("%s version %s by %s", params.MainParam, VERSION, AUTHOR)
	return result
}
