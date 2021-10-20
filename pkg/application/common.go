package application

import (
	"errors"
	"net/url"
	"strconv"
)

// --- AppError ---

// AppErrorType AppErrorのエラータイプ
type AppErrorType string

// AppByError エラーからのエラー
const AppByError AppErrorType = "error by error"
const NotFoundDataError AppErrorType = "not found data"

// String 文字列変換
func (errType AppErrorType) String() string {
	return string(errType)
}

// AppError UseCaseのエラー
type AppError interface {
	HasError() bool
	HaveType(AppErrorType) bool
	Error() string
	Type() AppErrorType
}

// UseCaseError AppErrorの実装
type UseCaseError struct {
	err     error
	errType AppErrorType
}

// NewError タイプだけを指定するコンストラクタ
func NewError(errType AppErrorType) AppError {
	return UseCaseError{
		err:     errors.New(errType.String()),
		errType: errType,
	}
}

// NewByError エラーを指定するコンストラクタ
func NewByError(err error) AppError {
	return UseCaseError{
		err:     err,
		errType: AppByError,
	}
}

// NewErrorWithMessage メッセージとタイプを指定するコンストラクタ
func NewErrorWithMessage(errType AppErrorType, message string) AppError {
	return UseCaseError{
		err:     errors.New(message),
		errType: errType,
	}
}

// HasError エラーがあるかどうか
func (err UseCaseError) HasError() bool {
	return err.err != nil
}

// HaveErrorType 指定のエラータイプを持っているか
func (err UseCaseError) HaveType(errType AppErrorType) bool {
	return err.HasError() && err.errType == errType
}

// Error エラーメッセージを返す
func (err UseCaseError) Error() string {
	return err.err.Error()
}

// Type エラータイプを返す
func (err UseCaseError) Type() AppErrorType {
	return err.errType
}

// --- Paginator ---

const (
	PaginationLimit int    = 10
	PaginationKey   string = "page"
)

// PaginatorElement ページングのためのページ番号に割り当てるURL
type PaginatorElement struct {
	Page int    `json:"page"`
	URL  string `json:"url"`
}

// Paginator ページング表示のためのDTO
type Paginator struct {
	TotalCount  int                `json:"total_count"`
	TotalPage   int                `json:"total_page"`
	Offset      int                `json:"offset"`
	Limit       int                `json:"limit"`
	CurrentPage int                `json:"current_page"`
	PrevPage    int                `json:"prev_page"`
	NextPage    int                `json:"next_page"`
	LastPage    int                `json:"last_page"`
	FirstCount  int                `json:"first_count"`
	LastCount   int                `json:"last_count"`
	PrevURL     string             `json:"prev_url"`
	NextURL     string             `json:"next_url"`
	Elements    []PaginatorElement `json:"elements"`
}

// SetURL ページングで利用するURLを設定する
func (p *Paginator) SetURL(original string) {
	u, _ := url.Parse(original)
	query := u.Query()

	// 前のページ番号のURL設定
	if p.CurrentPage > p.PrevPage {
		p.PrevURL = p.createPageNumberURL(p.PrevPage, u, query)
	} else {
		p.PrevURL = ""
	}

	// 次のページ番号のURL設定
	if p.CurrentPage < p.NextPage {
		p.NextURL = p.createPageNumberURL(p.NextPage, u, query)
	} else {
		p.NextURL = ""
	}

	// ページング表示の最初のページ番号
	firstPage := p.CurrentPage - PaginationLimit
	if 1 > firstPage {
		firstPage = 1
	}

	// ページング表示の最後のページ番号
	lastPage := p.CurrentPage + PaginationLimit
	if lastPage > p.TotalPage {
		lastPage = p.TotalPage
	}

	// ページング表示をする各ページ番号にURLを割り当てる
	elementLength := lastPage - (firstPage - 1)
	elements := make([]PaginatorElement, elementLength)

	count := 0
	for i := firstPage; i <= lastPage; i++ {
		elements[count] = PaginatorElement{
			Page: i,
			URL:  p.createPageNumberURL(i, u, query),
		}
		count++
	}
	p.Elements = elements
}

// createPageNumberURL ページ番号に対するURLの作成
func (p *Paginator) createPageNumberURL(
	page int,
	u *url.URL,
	q url.Values,
) string {
	q.Set(PaginationKey, strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return u.String()
}
