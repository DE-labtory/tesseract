package tesseract

type Tesseract struct {
}

// Deploy create Docker Container with running ShimCode and copying SmartContract.
func (t *Tesseract) Deploy() {
	// args : SmartContract info
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
