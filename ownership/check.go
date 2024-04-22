package ownership

import (
	"bytes"
	"fmt"
	"gomodules.xyz/go-sh"
	"k8s.io/klog/v2"
	"strings"
)

var (
	buf *bytes.Buffer
)

func ByExec() {
	shSession := getCommand("newer-0", "demo", 3)
	output, err := shSession.Output()
	if err != nil {
		klog.Infof("cant get output, err %s\n", err)
		return
	}

	errOutput := buf.String()
	if errOutput != "" {
		klog.Infof("failed to execute command, stderr: %s", errOutput)
		return
	}

	outStr := string(output)
	klog.Infoln(outStr + "\n\n")
	slice := strings.Split(outStr, "\n")
	for i, _ := range slice {
		if i == 0 {
			continue
		}
		klog.Infof("%s , ", slice[i])
	}
}

func checkOwnership() {

}

func getCommand(pod, ns string, column int) *sh.Session {
	shell := sh.NewSession()
	shell.ShowCMD = false
	buf = &bytes.Buffer{}
	shell.Stderr = buf

	podName := fmt.Sprintf("pod/%s", pod)
	kubectlCommand := []interface{}{
		"exec", "-n", ns, podName, "-c", "mongodb", "--",
	}
	mgCommand := []interface{}{
		"ls", "-l", "/data/db",
	}
	finalCommand := append(kubectlCommand, mgCommand...)

	awkCommand := fmt.Sprintf("{print $%d}", column)
	return shell.Command("kubectl", finalCommand...).Command("awk", awkCommand)
}
