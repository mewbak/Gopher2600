package regression

import (
	"fmt"
	"gopher2600/cartridgeloader"
	"path/filepath"
	"time"
)

// create a unique filename from a CatridgeLoader instance

func uniqueFilename(cartload cartridgeloader.Loader) string {
	n := time.Now()
	timestamp := fmt.Sprintf("%04d%02d%02d_%02d%02d%02d", n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second())
	newScript := fmt.Sprintf("%s_%s", cartload.ShortName(), timestamp)
	return filepath.Join(regressionScripts, newScript)
}
