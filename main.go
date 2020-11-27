package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	EnvVarConfigFile  = "CONFIG_FILE"
	DefaultConfigFile = "/tmp/config.yaml"
)

var (
	configFile = DefaultConfigFile
	mainLogger = ctrl.Log.WithName("pod-mutating-webhook")
)

func init() {
	if val, ok := os.LookupEnv(EnvVarConfigFile); ok {
		configFile = val
	}
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(false)))
	http.HandleFunc("/mutate", handleMutate)
	if err := http.ListenAndServeTLS(":8443", "/etc/webhook/certs/tls.crt", "/etc/webhook/certs/tls.key", nil); err != nil {
		mainLogger.Error(err, "Failed to start mutating webhook server")
	}
}

func handleMutate(w http.ResponseWriter, r *http.Request) {
	// load rules
	fileContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		mainLogger.Error(err, "Failed to read image rule config file")
		return
	}
	rules := &(map[string]string{})
	if err := yaml.Unmarshal(fileContents, rules); err != nil {
		mainLogger.Error(err, "Failed to unmarshall image rules.. skipping mutate")
		return
	}

	// read req body
	reqBody, errReading := ioutil.ReadAll(r.Body)
	if errReading != nil {
		mainLogger.Error(errReading, "Failed to read admission review req body")
		return
	}
	defer r.Body.Close()

	// mutate and get resp
	mutateImg := MutateContainerImage{logger: mainLogger.WithName("pod-image-mutator")}
	resp, errMutating := mutateImg.MutateContainerImages(reqBody, *rules)
	if errMutating != nil {
		mainLogger.Error(err, "Failed to mutate container image(s)")
		return
	}

	// meaning we have patches to apply
	if resp != nil {
		// write response back to k8s api
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		mainLogger.Info("Successfully mutated pod image(s)")
	}
}
