package tesseract

type Tesseract struct {
}

type ICodeInfo struct {
	Name      string
	Directory string
}

func (t *Tesseract) NewTesseract() {
}

// Deploy create Docker Container with running ShimCode and copying SmartContract.
func (t *Tesseract) SetupContainer() {
	// args : SmartContract info
	// Docker IMAGE pull
	// Create Docker
	// Copy ShimCode
	// Copy SmartContract
	// (Set DB info)
	// Running ShimCode on Container
	// (Connect socket)
	// Get Container handler
}

func (t *Tesseract) QueryOrInvoke() {
	// args : Transaction
	// Get Container handler using SmartContract ID
	// Send Query or Invoke massage
	// Receive result
	// Return result
}
