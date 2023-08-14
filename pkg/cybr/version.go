package cybr

import "fmt"

// Version field is a SemVer that should indicate the baked-in version of conceal
var Version = "1.0.0"

// Tag field denotes the specific build type for the broker. It may be replaced by compile-time variables if needed to
// provide the git commit information in the final binary.
var Tag = "release"

// FullVersionName is the user-visible aggregation of version and tag of this codebase
var FullVersionName = fmt.Sprintf("%s-%s", Version, Tag)
