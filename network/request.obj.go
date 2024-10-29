package network

type AddBlockRequest struct {
	Data string `json:"data"`
}

type AddNewPeerRequest struct {
	Address string `json:"address"`
}
