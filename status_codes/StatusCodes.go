package status_codes

const (
	RestartMarkerReply                = 110
	ServiceReadyInMinutes             = 120
	DataConnectionAlreadyOpen         = 125
	FileStatusOkay                    = 150
	CommandOkay                       = 200
	SystemStatus                      = 211
	DirectoryStatus                   = 212
	FileStatus                        = 213
	HelpMessage                       = 214
	NameSystemType                    = 215
	ServiceReadyForNewUser            = 220
	ServiceClosingControlConnection   = 221
	DataConnectionOpenNoTransfer      = 225
	ClosingDataConnection             = 226
	EnteringPassiveMode               = 227
	UserLoggedInProceed               = 230
	RequestedFileActionOkayCompleted  = 250
	PathnameCreated                   = 257
	UserNameOkayNeedPassword          = 331
	NeedAccountForLogin               = 332
	RequestedFileActionPending        = 350
	ServiceNotAvailable               = 421
	CantOpenDataConnection            = 425
	ConnectionClosedTransferAborted   = 426
	RequestedFileActionNotTaken       = 450
	RequestedActionAborted            = 451
	SyntaxErrorCommandUnrecognized    = 500
	SyntaxErrorParametersArguments    = 501
	CommandNotImplemented             = 502
	BadSequenceOfCommands             = 503
	CommandNotImplementedForParameter = 504
	NotLoggedIn                       = 530
	NeedAccountForStoringFiles        = 532
	RequestedActionNotTaken           = 550
	RequestedActionAbortedPageType    = 551
	RequestedFileActionAborted        = 552
	RequestedActionNotTakenFileName   = 553
	IntegrityProtectedReply           = 631
	ConfidentialityIntegrityProtected = 632
	ConfidentialityProtectedReply     = 633
)
