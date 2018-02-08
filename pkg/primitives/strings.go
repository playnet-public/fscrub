package primitives

// BuildHeader returns []string containing information about fscrub
func BuildHeader() []string {
	return []string{
		"//// PlayNet Fscrub ////",
		"// This file got cleaned by fscrub (github.com/playnet-public/fscrub) to remove sensitive information.",
		"// If this action was taken by mistake, please contact your responsible admin or seek advice at PlayNet (https://discord.gg/vhbP6Ks).",
		"// To make fscrub ignore your file, please add \"//-ignore: github.com/playnet-public/fscrub\" as first line and upload it again.",
		"////",
	}
}

// BuildIgnoreHeader returns the header required to make fscrub ignore a file
func BuildIgnoreHeader() string {
	return "//-ignore: github.com/playnet-public/fscrub"
}
