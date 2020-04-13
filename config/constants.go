package config

const (
	LogPath = "./Log"
)

const (
	NASA_NAME    = "nasa"
	NASA_API_KEY = "FnRwnqD0SGyf2GUbEwaAH33H4d6TeaRYYfdOEZwl"
)

const (
	ONT_MAIN_NET    = "http://dappnode1.ont.io:20336"
	ONT_TEST_NET    = "http://polaris1.ont.io:20336"
	ONT_SOLO_NET    = "http://127.0.0.1:20336"
	VERIFY_TX_RETRY = 6
)

const (
	NETWORK_ID_MAIN_NET    = 1
	NETWORK_ID_POLARIS_NET = 2
	NETWORK_ID_SOLO_NET    = 3
)

type TxStatus uint8

const (
	Paying TxStatus = iota
	PayFailed
	PaySuccess
)
