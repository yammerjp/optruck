package exec

import execPackage "k8s.io/utils/exec"

type ExecInterface execPackage.Interface
type ExecCommand execPackage.Cmd

var exec execPackage.Interface

func init() {
	exec = execPackage.New()
}

func GetExec() execPackage.Interface {
	return exec
}

func SetExec(e execPackage.Interface) {
	exec = e
}
