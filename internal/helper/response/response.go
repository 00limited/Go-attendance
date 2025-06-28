package response

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type responseHelper struct {
	Validation   interface{}
	ViewNilValid bool
}

type responseFormat struct {
	C       echo.Context
	Code    int
	Status  string
	Message string
	Data    interface{}
}

// Interface ...
type Interface interface {
	SetResponse(c echo.Context, code int, status string, message string, data interface{}) responseFormat
	SendResponse(res responseFormat) error
	EmptyJSONMap() map[string]interface{}
	SendSuccess(c echo.Context, message string, data interface{}) error
	SendSuccessWithValidation(c echo.Context, message string, data interface{}, validation interface{}, viewNilvalid bool) error
	SendBadRequest(c echo.Context, message string, data interface{}) error
	SendError(c echo.Context, message string, data interface{}) error
	SendErrorWithValidation(c echo.Context, message string, data interface{}, validation interface{}) error
	SendUnauthorized(c echo.Context, message string, data interface{}) error
	SendValidationError(c echo.Context, validationErrors validator.ValidationErrors) error
	SendNotFound(c echo.Context, message string, data interface{}) error
	SendCustomResponse(c echo.Context, httpCode int, message string, data interface{}) error
	SendResponsByCode(c echo.Context, code int, message string, data interface{}, err error) error
	SendPaginationResponse(c echo.Context, items interface{}, message string, totalRecord, totalRecordPerPage, totalRecordSearch, totalPage int64, currentPage int) error
	GetBranch() string
	GetHash(branch string) string
	GetUpdated() string
	GetHostname() string
}

// NewResponse ...
func NewResponse() Interface {
	return &responseHelper{}
}

// SetResponse ...
func (r *responseHelper) SetResponse(c echo.Context, code int, status string, message string, data interface{}) responseFormat {

	return responseFormat{c, code, status, message, data}
}

// SendResponse ...
func (r *responseHelper) SendResponse(res responseFormat) error {
	if len(res.Message) == 0 {
		res.Message = http.StatusText(res.Code)
	}
	resp := map[string]interface{}{
		// "git_branch":  r.GetBranch(),
		// "git_hash":    r.GetHash(r.GetBranch()),
		// "git_updated": r.GetUpdated(),
		// "hostname":    r.GetHostname(),
		"code":    res.Code,
		"status":  res.Status,
		"message": res.Message,
	}
	if res.Data != nil {
		resp["data"] = res.Data
	}

	// fmt.Printf("r.Validation: %s\n", reflect.TypeOf(r.Validation))
	// // fmt.Println("r.Validation != nil ?", (r.Validation != nil))
	// var dataValid []string
	// typeValid := reflect.TypeOf(r.Validation)
	// if typeValid.String() == reflect.TypeOf(dataValid).String() {
	// 	if r.Validation != nil {
	// 		dataValid = r.Validation.([]string)
	// 	}
	// }
	if r.Validation != nil {
		resp["validation"] = r.Validation
		r.Validation = nil
		println("after r.Validation : ", r.Validation)
	} else if r.ViewNilValid {
		resp["validation"] = nil
		r.ViewNilValid = false
	}
	return res.C.JSON(res.Code, resp)
}

func (r *responseHelper) GetBranch() string {
	branchName := ""

	firstLineHead := ""
	file, err := os.Open(".git/HEAD")
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		firstLineHead = scanner.Text()
		break
	}
	headInfo := strings.Split(firstLineHead, "/")
	if len(headInfo) >= 2 {
		for i := 2; i < len(headInfo); i++ {
			if i > 2 {
				branchName += "/"
			}
			branchName += headInfo[i]
		}
	}

	return branchName
}

func (r *responseHelper) GetHash(branchName string) string {
	hash := ""

	file, err := os.Open(".git/refs/heads/" + branchName)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hash = scanner.Text()
		break
	}

	return hash
}

func (r *responseHelper) GetUpdated() string {
	statinfo, err := os.Stat(".git/index")
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return statinfo.ModTime().Format("2006-01-02 15:04:05")
}

func (r *responseHelper) GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "-"
	}

	return hostname
}

func (r *responseHelper) SendResponsByCode(c echo.Context, code int, message string, data interface{}, err error) error {
	if err != nil {
		message = err.Error()
	}

	res := r.SetResponse(c, code, http.StatusText(code), message, data)
	return r.SendResponse(res)
}

// EmptyJSONMap : set empty data.
func (r *responseHelper) EmptyJSONMap() map[string]interface{} {
	return make(map[string]interface{})
}

// SendSuccess : Send success response to consumers.
func (r *responseHelper) SendSuccess(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusOK, http.StatusText(http.StatusOK), message, data)
	return r.SendResponse(res)
}

// SendBadRequest : Send bad request response to consumers.
func (r *responseHelper) SendBadRequest(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), message, data)
	return r.SendResponse(res)
}

// SendError : Send error request response to consumers.
func (r *responseHelper) SendError(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), message, data)
	return r.SendResponse(res)
}

// SendUnauthorized : Send error request response to consumers.
func (r *responseHelper) SendUnauthorized(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), message, data)
	return r.SendResponse(res)
}

// SendValidationError : Send validation error request response to consumers.
func (r *responseHelper) SendValidationError(c echo.Context, validationErrors validator.ValidationErrors) error {
	errorResponse := []string{}
	for _, err := range validationErrors {
		errorResponse = append(errorResponse, strings.Trim(fmt.Sprint(err), "[]")+".")
	}
	res := r.SetResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), strings.Trim(fmt.Sprint(errorResponse), "[]"), r.EmptyJSONMap())
	return r.SendResponse(res)
}
func (r *responseHelper) SendErrorWithValidation(c echo.Context, message string, data interface{}, validation interface{}) error {
	r.Validation = validation
	return r.SendError(c, message, data)
}
func (r *responseHelper) SendSuccessWithValidation(c echo.Context, message string, data interface{}, validation interface{}, viewNilValid bool) error {
	r.Validation = validation
	r.ViewNilValid = viewNilValid
	return r.SendSuccess(c, message, data)
}

// SendNotFound : Send error request response to consumers.
func (r *responseHelper) SendNotFound(c echo.Context, message string, data interface{}) error {
	res := r.SetResponse(c, http.StatusNotFound, http.StatusText(http.StatusNotFound), message, data)
	return r.SendResponse(res)
}

// SendCustomResponse ...
func (r *responseHelper) SendCustomResponse(c echo.Context, httpCode int, message string, data interface{}) error {
	res := r.SetResponse(c, httpCode, http.StatusText(httpCode), message, data)
	return r.SendResponse(res)
}

// Send Pagination Response
type ResponsePagination struct {
	Records            interface{} `json:"records"`
	TotalRecord        int64       `json:"total_record"`
	TotalRecordPerPage int64       `json:"total_record_per_page"`
	TotalRecordSearch  int64       `json:"total_record_search"`
	TotalPage          int64       `json:"total_page"`
	CurrentPage        int         `json:"current_page"`
	RowNumberStart     int         `json:"row_number_start"`
	RowNumberEnd       int         `json:"row_number_end"`
}

func (r *responseHelper) SendPaginationResponse(c echo.Context, items interface{}, message string, totalRecord, totalRecordPerPage, totalRecordSearch, totalPage int64, currentPage int) error {
	rowNumberStart := (currentPage-1)*int(totalRecordPerPage) + 1
	var rowNumberEnd int
	if currentPage == int(totalPage) {
		rowNumberEnd = int(totalRecord)
	} else {
		rowNumberEnd = currentPage * int(totalRecordPerPage)
	}

	response := &ResponsePagination{
		Records:            items,
		TotalRecord:        totalRecord,
		TotalRecordPerPage: totalRecordPerPage,
		TotalRecordSearch:  totalRecordSearch,
		TotalPage:          totalPage,
		CurrentPage:        currentPage,
		RowNumberStart:     rowNumberStart,
		RowNumberEnd:       rowNumberEnd,
	}

	res := r.SetResponse(c, http.StatusOK, http.StatusText(http.StatusOK), message, response)

	return r.SendResponse(res)
}
