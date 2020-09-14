package mtest

import . "github.com/onsi/ginkgo"

// FunctionsSuite is a test suite that tests small test cases
var FunctionsSuite = func() {
	Context("coil-controller", TestCoilController)
	Context("coil-installer", TestCoilInstaller)
	Context("pod", TestPod)
	Context("pool", TestPool)
}

// FailuresSuite is a test suite that runs test cases with failure injection
var FailuresSuite = func() {
	Context("pod startup", TestPodStartup)
}

// BootstrapSuite is a test suite that bootstrap demo environment
var BootstrapSuite = func() {
	Context("setup-coil", func() {
		It("should setup coil", func() {
			initializeCoil()
		})
	})
}
