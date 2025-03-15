package semver

func Compare(v1, v2 Version) int { return v1.Compare(v2) }

func Less(v1, v2 Version) bool { return v1.LessThan(v2) }
