package parameter

import "strings"

const CLI_TRUE_STR = "true"
const CLI_FALSE_STR = "false"

// Parameters represents a structured way to handle CLI parameters, including both key-value pairs
// and a main parameter for commands that require a primary argument.
type Parameters struct {
	// Program holds the name of the CLI program being executed.
	Program string
	// Params stores key-value pairs of CLI options and their associated values.
	Params map[string]string
	// MainParam captures the primary argument for commands that require a main parameter
	// (e.g., a username for add-user).
	MainParam string
}

// Parse converts CLI arguments into a structured Parameters instance, separating key-value
// pairs and the main parameter.
func Parse(cmds ...string) *Parameters {
	mainParam := ""
	program := ""
	params := make(map[string]string)

	if len(cmds) > 0 {
		program = cmds[0]
	}

	if len(cmds) > 1 {
		cmds = cmds[1:]

		for i := 0; i < len(cmds); i++ {
			param := cmds[i]

			if isKeyword(param) {
				key := strings.TrimLeft(param, "-")

				if i+1 < len(cmds) && !isKeyword(cmds[i+1]) {
					params[key] = cmds[i+1]
					i++
				} else {
					params[key] = CLI_TRUE_STR
				}

				continue
			}

			if mainParam == "" {
				mainParam = param
			}
		}
	}

	return &Parameters{
		Params:    params,
		MainParam: mainParam,
		Program:   program,
	}
}

// isKeyword checks if the provided string is a CLI option (starts with "-" or "--").
func isKeyword(s string) bool {
	s = strings.ToLower(s)
	if (strings.HasPrefix(s, "--") || strings.HasPrefix(s, "-")) && !strings.HasPrefix(s, "---") {
		return true
	}

	return false
}
