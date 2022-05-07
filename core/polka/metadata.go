package polka

// Load cached metadata string.
// This will save a lot of network traffic to download metadata from rpcUrl.
func (c *Chain) LoadCachedMetadataString(metadataString string) error {
	_, err := getOrCreatePolkaClient(c.RpcUrl, metadataString)
	return err
}

// Get the metadata string of the chain (if not, it will be downloaded automatically)
func (c *Chain) GetMetadataString() (string, error) {
	client, err := getConnectedPolkaClient(c.RpcUrl)
	if err != nil {
		return "", nil
	}
	return client.MetadataString()
}

// Reload the latest metadata of this chain.
// @return the latest metadata string
func (c *Chain) ReloadMetadata() (string, error) {
	client, err := getConnectedPolkaClient(c.RpcUrl)
	if err != nil {
		return "", nil
	}

	err = client.ReloadMetadata()
	if err != nil {
		return "", nil
	}

	return client.MetadataString()
}
