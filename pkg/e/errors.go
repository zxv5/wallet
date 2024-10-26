package e

var (
	OK = add(0) // OK .

	// default server

	NotModified           = add(304).US("Not Modified")
	RequestErr            = add(400).US("Request Error")
	Unauthorized          = add(401).US("Unauthorized")
	AccessDenied          = add(403).US("Access Denied")
	NotFound              = add(404).US("Not Found")
	MethodNotAllowed      = add(405).US("Method Not Allowed")
	Conflict              = add(409).US("Conflict")
	Canceled              = add(498).US("context canceled")
	ServerErr             = add(500).US("Server Error")
	ServiceUnavailable    = add(503).US("Service Unavailable")
	Deadline              = add(504).US("Deadline")
	LimitExceed           = add(509).US("Limit Exceed")
	FileNotExists         = add(616).US("File Not Exists")
	FileTooLarge          = add(617).US("File Too Large")
	FailedTooManyTimes    = add(625).US("Failed Too Many Times")
	UserNotExist          = add(626).US("User Not Exist")
	PasswordTooLeak       = add(628).US("Password Too Leak")
	UsernameOrPasswordErr = add(629).US("Username Or Password Error")
	TargetNumberLimit     = add(632).US("Target Number Limit")
	TargetBlocked         = add(643).US("Target Blocked")
	UserLevelLow          = add(650).US("User Level Too Low")
	UserDuplicate         = add(652).US("User Duplicate")
	AccessTokenExpires    = add(658).US("AccessToken Expires")
	PasswordHashExpires   = add(662).US("Password Expires")
	UnexpectedErr         = add(701).US("Unexpected Error")
	SQLErr                = add(710).US("Sql Error")
	RecordNotFoundErr     = add(711).US("Not Found Data")
	DuplicatedErr         = add(712).US("Duplicated")

	// common error
	AuthErr          = add(10001).US("Authentication Error")
	PermissionDenied = add(10002).US("Permission denied")
	RequestException = add(10003).US("Request Exception")

	ParamsErr = add(11001).US("Params Error")
	CheckErr  = add(11002).US("Validate Error")

	// error
	UserNotActive     = add(30002).US("The account has not been activated.")
	InsufficientFunds = add(30003).US("Insufficient funds")
)
