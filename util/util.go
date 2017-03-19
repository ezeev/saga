// Package util provides "glue" to the cloud-ninja packages
package util

import (
	"os"
	"fmt"
)


// checkVars ensures that required environment variables are set. Pass a slice of required env var keys.
func CheckVars(vars []string) error {
	missingVars := ""
	for _, k := range vars {
		if os.Getenv(k) == "" {
			missingVars = fmt.Sprintf(missingVars + ", %s", k)
		}
	}
	if missingVars != "" {
		return fmt.Errorf("Required Environment variables are missing: %s", missingVars)
	} else {
		return nil
	}
}

