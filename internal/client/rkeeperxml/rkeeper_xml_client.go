package rkeeperxml

type RKeeperXMLClient struct{}

func New() *RKeeperXMLClient {
	return &RKeeperXMLClient{}
}

func (c *RKeeperXMLClient) GetOrderList() error {
	return nil
}
