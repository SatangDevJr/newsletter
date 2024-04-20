package error

import (
	"net/http"

	"github.com/thoas/go-funk"
)

type ErrorCode string

type ErrorMessage struct {
	EN         string
	TH         string
	StatusCode int
	Code       ErrorCode
}

const (
	TH = "th"
	EN = "en"
)

const (
	InternalServerErrorMessage = "Internal server error"
	DataNotFoundMessage        = "Data not found"
	BadRequestMessage          = "Bad request"
	ConflictMessage            = "Conflict"
	UnauthorizedMessage        = "Unauthorized"
	ForbiddenMessage           = "Forbidden"
	DataQuotaNotFoundMessage   = "Data quota not found"
	IncorrectDataMessage       = "Incorrect Data"
)

// error code
const (
	InternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	DataNotFound        ErrorCode = "DATA_NOT_FOUND"
	BadRequest          ErrorCode = "BAD_REQUEST"
	Conflict            ErrorCode = "CONFLICT"
	Unauthorized        ErrorCode = "UNAUTHORIZED"
	Forbidden           ErrorCode = "FORBIDDEN"

	DatabaseConnectionFailed ErrorCode = "DATABASE_CONNECTION_FAILED"
	TechnicalError           ErrorCode = "TECHNICAL_ERROR"

	InvalidLanguage          ErrorCode = "INVALID_LANGUAGE"
	ValidateQuotaNotEnough   ErrorCode = "VALIDATE_QUOTA_NOTENOUGH"
	DuplicateFileName        ErrorCode = "DUPLICATE_FILE_NAME"
	DuplicateCardId          ErrorCode = "DUPLICATE_CARD_ID"
	DuplicateTaxRegistration ErrorCode = "DUPLICATE_TAX_REGISTRATION"
	DuplicateUserName        ErrorCode = "DUPLICATE_USER_NAME"
	InvalidPassword          ErrorCode = "INVALID_PASSWORD"
	DataQuotaNotFound        ErrorCode = "DATA_QUOTA_NOT_FOUND"
	EmptyToken               ErrorCode = "EMPTY_TOKEN"
	IncorrectUserName        ErrorCode = "INCORRECT_USERNAME"
	QuotaAlreadyUsed         ErrorCode = "QUOTA_ALREADY_USED"
	AmountUseMoreThanZero    ErrorCode = "AmountUse_More_than_Zero"
	CanNotDeleteInvoice      ErrorCode = "CAN_NOT_DELETE_INVOICE"
	CanNotDeleteDeliveryNote ErrorCode = "CAN_NOT_DELETE_DELIVERY_NOTE"
	UploadFailed             ErrorCode = "UPLOAD_FAILED"
	ExistUserName            ErrorCode = "EXIST_USER_NAME"
	UserActiveFailed         ErrorCode = "USER_ACTIVE_FAILED"
	SubPartnerHasAnInvoice   ErrorCode = "SUB_PARTNER_HAS_AN_INVOICE"
)

var mappingMessage = map[ErrorCode]ErrorMessage{
	InternalServerError: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         InternalServerErrorMessage,
		TH:         "ระบบขัดข้อง",
	},
	DataNotFound: {
		Code:       DataNotFound,
		StatusCode: http.StatusNotFound,
		EN:         DataNotFoundMessage,
		TH:         "ไม่พบข้อมูล",
	},
	DataQuotaNotFound: {
		Code:       DataQuotaNotFound,
		StatusCode: http.StatusNotFound,
		EN:         DataQuotaNotFoundMessage,
		TH:         "ไม่พบข้อมูลโควต้า",
	},
	BadRequest: {
		Code:       BadRequest,
		StatusCode: http.StatusBadRequest,
		EN:         BadRequestMessage,
		TH:         "คำขอไม่ถูกต้อง",
	},
	Conflict: {
		Code:       Conflict,
		StatusCode: http.StatusConflict,
		EN:         ConflictMessage,
		TH:         "ข้อมูลขัดแย้งในระบบ",
	},
	Unauthorized: {
		Code:       Unauthorized,
		StatusCode: http.StatusUnauthorized,
		EN:         UnauthorizedMessage,
		TH:         "สิทธิ์การเข้าใช้งานหมดอายุ",
	},
	Forbidden: {
		Code:       Forbidden,
		StatusCode: http.StatusForbidden,
		EN:         ForbiddenMessage,
		TH:         "ไม่มีสิทธิ์เข้าถึง",
	},
	DatabaseConnectionFailed: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "Database connection failed.",
		TH:         "ไม่สามารถเชื่อมต่อฐานข้อมูล",
	},
	InvalidLanguage: {
		Code:       BadRequest,
		StatusCode: http.StatusBadRequest,
		EN:         "Invalid accpet language",
		TH:         "ภาษาที่ส่งมาไม่ถูกต้อง",
	},
	ValidateQuotaNotEnough: {
		Code:       BadRequest,
		StatusCode: http.StatusBadRequest,
		EN:         "Unable to add quota because quota of month not enough",
		TH:         "ไม่สามารถเพิ่มจำนวน Quota ที่ท่านระบุได้เนื่องจากมีการใช้งานครบตามจำนวน Quota ที่ได้รับภายในเดือน",
	},
	DuplicateFileName: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "Duplicate file name",
		TH:         "ชื่อไฟล์ซ้ำ กรุณาตั้งชื่อไฟล์ใหม่",
	},
	DuplicateCardId: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "Duplicate card id",
		TH:         "บัตรประชาชนซ้ำ กรุณากรอกใหม่อีกครั้ง",
	},
	DuplicateTaxRegistration: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "Duplicate tax registration",
		TH:         "ใบกำกับภาษีซ้ำ กรุณากรอกใหม่อีกครั้ง",
	},
	DuplicateUserName: {
		Code:       Conflict,
		StatusCode: http.StatusConflict,
		EN:         "Duplicate user name",
		TH:         "รหัสผู้ใช้ซ้ำ กรุณากรอกใหม่อีกครั้ง",
	},
	ExistUserName: {
		Code:       ExistUserName,
		StatusCode: http.StatusConflict,
		EN:         "Existed user name",
		TH:         "มีชื่อผู้ใช้งานนี้ในระบบแล้ว",
	},
	InvalidPassword: {
		Code:       BadRequest,
		StatusCode: http.StatusBadRequest,
		EN:         "ชื่อผู้ใช้งาน หรือ รหัสผ่านไม่ถูกต้อง",
		TH:         "ชื่อผู้ใช้งาน หรือ รหัสผ่านไม่ถูกต้อง",
	},
	EmptyToken: {
		Code:       BadRequest,
		StatusCode: http.StatusBadRequest,
		EN:         "Empty Token",
		TH:         "Empty Token",
	},
	IncorrectUserName: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "Incorrect username",
		TH:         "รหัส username ไม่ถูกต้อง",
	},
	QuotaAlreadyUsed: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "Can't update quota. It's already been used.",
		TH:         "โควต้าถูกใช้งานไปแล้วไม่สามารถแก้ไขได้",
	},
	AmountUseMoreThanZero: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "AmountUseMorethanZero.",
		TH:         "โควต้าถูกใช้งานไปแล้วไม่สามารถแก้ไขได้",
	},
	CanNotDeleteInvoice: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "Can't delete invoice.",
		TH:         "ไม่สามารถลบเอกสารวางบิลได้",
	},
	CanNotDeleteDeliveryNote: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "Can't delete delivery note.",
		TH:         "ไม่สามารถลบใบส่งสินค้าได้",
	},
	UploadFailed: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "Failed to upload file",
		TH:         "อัพโหลดไฟล์ไม่สําเร็จ",
	},
	UserActiveFailed: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "User ของท่านได้ถูก Inactive กรุณาติดต่อ ADMIN เพื่อเปิดใช้งาน",
		TH:         "User ของท่านได้ถูก Inactive กรุณาติดต่อ ADMIN เพื่อเปิดใช้งาน",
	},
	SubPartnerHasAnInvoice: {
		Code:       InternalServerError,
		StatusCode: http.StatusInternalServerError,
		EN:         "ไม่สามารถทำการลบช้อมูลลูกข่ายได้ เนื่องจากมีรายการคงค้างในระบบ",
		TH:         "ไม่สามารถทำการลบช้อมูลลูกข่ายได้ เนื่องจากมีรายการคงค้างในระบบ",
	},
}

func MapMessageError(code ErrorCode, languageCode string) (statusCode int, e Error) {

	if _, ok := mappingMessage[code]; !ok || !funk.ContainsString([]string{TH, EN}, languageCode) {
		e.Code = InternalServerError
		statusCode = http.StatusInternalServerError
		if languageCode == TH {
			e.Message = "เกิดข้อผิดพลาดบางอย่าง"
		} else {
			e.Message = "Something was wrong."
		}
		return
	}

	statusCode = mappingMessage[code].StatusCode
	e.Code = mappingMessage[code].Code
	switch language := languageCode; language {
	case TH:
		e.Message = mappingMessage[code].TH
	default:
		e.Message = mappingMessage[code].EN
	}

	return
}
