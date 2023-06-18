package nix

import "os"

const nixDebugEnv string = "NIXFAAS_NIX_DEBUG"

var nixDebug bool = getEnv(nixDebugEnv, "false") == "true"

func getEnv(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return v
}
