package sagaconfig

const (
	LogPath = "./Log"
)

const (
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
	NETWORK_ID_TRAVIS_NET  = 4
)

const (
	CategoryAll = 1
)

type OrderStatus uint8

const (
	Processing OrderStatus = iota
	Canceled
	Failed
	Completed
)

const (
	ONG_DECIMALS = 9
)

const (
	Api  = "api"
	Data = "data"
)

//api token type
const (
	TOKEN_TYPE_ONG        = "ONG"
	TOKEN_TYPE_ONT        = "ONT"
	ONG_CONTRACT_ADDRESS  = "0200000000000000000000000000000000000000"
	Collect_Money_Address = "AbtTQJYKfQxq4UdygDsbLVjE8uRrJ2H3tP"
)

const (
	QrCodeExp       = 10 * 60
	DefRequestLimit = 100
)

const (
	QueryAmt     = 3
	Key_OntId    = "OntId"
	OntId        = "did:ont:AYCcjQuB6xgXm2vKku9Vb6bdTcEguXqbt1"
	OntIdPrivate = "ae1bab4364ec7966ab8e8a1db43cf7162b6e619bcab9ce0af4d1763bc4a62186"
)
const (
	APOD = 1
	FEED = 2
)
