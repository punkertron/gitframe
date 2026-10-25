package cli

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

func printProgress(cntFiles int, statusWriter io.Writer, done <-chan struct{}, cnt *atomic.Int64) {
	if cntFiles <= 0 {
		panic(fmt.Sprintf("invalid file count: %d", cntFiles))
	}

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	printed := false
	defer func() {
		if printed {
			fmt.Fprint(statusWriter, "\r\x1b[2K")
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			amount := cnt.Load()
			fmt.Fprintf(statusWriter, "\rProgress: %3d%% (analyzed %d files out of %d)", amount*100/int64(cntFiles), amount, cntFiles)
			printed = true
		}
	}
}
