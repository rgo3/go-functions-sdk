package deploy

import (
	"fmt"
	"go/types"
	"os/exec"
	"path"
	"plugin"
	"strconv"
	"strings"

	"github.com/dergoegge/go-functions-sdk/internal/pkg/build"
	"github.com/dergoegge/go-functions-sdk/internal/pkg/config"
	"github.com/dergoegge/go-functions-sdk/internal/pkg/parse"
	"github.com/dergoegge/go-functions-sdk/pkg/functions"
)

var projectID string

func bucketNameToResource(bucket string) string {
	if strings.Compare(bucket, "default") == 0 {
		bucket = projectID
	}

	return fmt.Sprintf("gs://%s.appspot.com", bucket)
}

func firestoreDocToResource(doc string) string {
	return fmt.Sprintf("projects/%s/databases/(default)/documents/%s", projectID, doc)
}

func getDeployFlags(fnSym plugin.Symbol) ([]string, *functions.FunctionBuilder) {
	var triggerFlags []string
	var builder *functions.FunctionBuilder

	switch fnSym.(type) {
	case **functions.HTTPFunctionBuilder:
		httpBuilder := (*fnSym.(**functions.HTTPFunctionBuilder))
		builder = httpBuilder.FunctionBuilder
		triggerFlags = []string{"--trigger-http"}
		break

	case **functions.FirestoreFunctionBuilder:
		firestoreBuilder := (*fnSym.(**functions.FirestoreFunctionBuilder))
		builder = firestoreBuilder.FunctionBuilder
		triggerFlags = []string{
			"--trigger-event", firestoreBuilder.Event,
			"--trigger-resource", firestoreDocToResource(firestoreBuilder.Doc),
		}
		break
	case **functions.StorageFunctionBuilder:
		storageBuilder := (*fnSym.(**functions.StorageFunctionBuilder))
		builder = storageBuilder.FunctionBuilder
		triggerFlags = []string{
			"--trigger-event", storageBuilder.Event,
			"--trigger-resource", bucketNameToResource(storageBuilder.GCBucket),
		}
		break
	case **functions.FirebaseAuthFunctionBuilder:
		authBuilder := (*fnSym.(**functions.FirebaseAuthFunctionBuilder))
		builder = authBuilder.FunctionBuilder
		triggerFlags = []string{
			"--trigger-event", authBuilder.Event,
			"--trigger-resource", projectID,
		}
		break
	}

	return triggerFlags, builder
}

func createDeployCommand(fnPlugin *plugin.Plugin, fn types.Object, pkgPath string) (*exec.Cmd, error) {
	fnSym, err := fnPlugin.Lookup(fn.Name())
	if err != nil {
		return nil, err
	}

	deployFlags, fnBuilder := getDeployFlags(fnSym)

	cmdArgs := []string{
		"functions",
		"deploy",
		fn.Name(),
		"--runtime", "go111",
		"--source", "./" + pkgPath,
		"--entry-point", fn.Name() + ".Handler",
		"--region", fnBuilder.GCRegion,
		"--memory", fnBuilder.RuntimeOpts.Memory,
		"--timeout", strconv.Itoa(fnBuilder.RuntimeOpts.Timeout),
	}
	cmdArgs = append(cmdArgs, deployFlags...)

	return exec.Command("gcloud", cmdArgs...), nil
}

func preparePackage(pkg parse.Package) ([]*exec.Cmd, error) {
	fnPlugin, err := plugin.Open(path.Join(build.PluginFolder, pkg.Name+".so"))
	if err != nil {
		return nil, err
	}

	cmds := []*exec.Cmd{}
	for _, fn := range pkg.Functions {
		cmd, err := createDeployCommand(fnPlugin, fn, pkg.Path)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

// Prepare creates deploy commands for the functions in the provided packages.
func Prepare(fnPackages parse.Packages) ([]*exec.Cmd, error) {
	conf, err := config.NewSDKConfig()
	if err != nil {
		return nil, err
	}

	projectID, err = conf.ProjectID()
	if err != nil {
		return nil, err
	}

	cmds := []*exec.Cmd{}
	for _, pkg := range fnPackages {
		pkgCmds, err := preparePackage(pkg)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, pkgCmds...)
	}

	return cmds, nil
}
