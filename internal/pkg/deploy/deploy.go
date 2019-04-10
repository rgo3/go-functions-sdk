package deploy

import (
	"bytes"
	"fmt"
	"go/types"
	"os/exec"

	"github.com/google/logger"
)

func deploy(fn types.Object, cmd *exec.Cmd, done chan error) {
	logger.Infof("Deploying %s", fn.Name())
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf
	err := cmd.Run()

	if err == nil {
		logger.Infof("Deployed %s successfully", fn.Name())
	} else {
		err = fmt.Errorf("Deployment of %s failed:\n%s", fn.Name(), outBuf.String())
	}

	done <- err
}

func Functions(stagedFunctions []types.Object, deployCmds []*exec.Cmd) error {
	wait := make(chan error)

	for i, fn := range stagedFunctions {
		go deploy(fn, deployCmds[i], wait)
	}

	for range stagedFunctions {
		err := <-wait
		if err != nil {
			logger.Error(err)
		}
	}
	return nil
}
