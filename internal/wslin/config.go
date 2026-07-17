package wslin

import (
	"os"
	"path/filepath"
)

var InstallRoot = filepath.Join(os.Getenv("LOCALAPPDATA"), "wslbm", "wslos")
