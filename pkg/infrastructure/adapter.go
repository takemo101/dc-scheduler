package infrastructure

import (
	"mime/multipart"
	"os"

	"github.com/takemo101/dc-scheduler/core"
	"github.com/valyala/fasthttp"
)

// --- UploadAdapter ---

// UploadAdapter ファイルアップロード機能
type UploadAdapter interface {
	Upload(file *multipart.FileHeader, path string) (string, error)
	Delete(path string) error
	MakeDirectory(directory string) error
	Exists(path string) bool
	ToPath(path string) string
	ToURL(path string) string
}

// PublicStorageUploadAdapter 公開領域にアップロードするUploadAdapterの実装
type PublicStorageUploadAdapter struct {
	path core.Path
}

// NewPublicStorageUploadAdapter コンストラクタ
func NewPublicStorageUploadAdapter(
	path core.Path,
) UploadAdapter {
	return PublicStorageUploadAdapter{
		path,
	}
}

// Upload アップロードを実行してアップロード先のフルパスを返す
func (ap PublicStorageUploadAdapter) Upload(file *multipart.FileHeader, path string) (toPath string, err error) {

	toPath = ap.ToPath(path)

	err = fasthttp.SaveMultipartFile(file, toPath)

	if err != nil {
		return toPath, err
	}

	return toPath, err
}

// Delete 相対パスからファイルを削除する
func (ap PublicStorageUploadAdapter) Delete(path string) (err error) {
	return os.Remove(
		ap.ToPath(path),
	)
}

// MakeDirectory 相対ディレクトリパスからディレクトリを作成する
func (ap PublicStorageUploadAdapter) MakeDirectory(directory string) (err error) {

	dir := ap.ToPath(directory)
	if f, e := os.Stat(dir); os.IsNotExist(e) || !f.IsDir() {
		return os.MkdirAll(dir, 0777)
	}

	return err
}

// Exists ファイルディレクトリの存在するか
func (ap PublicStorageUploadAdapter) Exists(path string) bool {

	_, e := os.Stat(path)
	return !os.IsNotExist(e)
}

// ToPath 相対パスからフルパスを返す
func (ap PublicStorageUploadAdapter) ToPath(path string) (toPath string) {
	return ap.path.Public(path)
}

// ToURL 相対パスから公開URLを返す
func (ap PublicStorageUploadAdapter) ToURL(path string) (toPath string) {
	return ap.path.PublicURL(path)
}
