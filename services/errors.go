package services

var (
	MsgCodeOK                        = 0
	MsgCodeUnauthorized              = 2
	MsgCodeDBUnexpectedErr           = 100
	MsgCodeDBRecordsNotFound         = 101
	MsgCodeHelperCurrentUserNotFound = 200
	MsgCodeReqHelperNotJSON          = 1300
	MsgCodeReqHelperBadlyFormedAtPos = 1301
	MsgCodeReqHelperBadlyFormed      = 1302
	MsgCodeReqHelperInvalidValue     = 1303
	MsgCodeReqHelperUnknownField     = 1304
	MsgCodeReqHelperReqBodyEmpty     = 1305
	MsgCodeReqHelperLimitSize        = 1306
	MsgCodeReqHelperLimit1Obj        = 1307
	MsgCodeProcessOkWithErrs         = 1400
)
