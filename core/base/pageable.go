package base

type Pageable interface {
	Jsonable

	// The total count of all data in the remote server. Returns 0 if statistics are not available
	TotalCount() int
	// The count of data in the current page.
	CurrentCount() int
	// The cursor of the current page.
	CurrentCursor() string
	// Is there has next page.
	HasNextPage() bool

	ItemArray() *AnyArray

	// You need to implement the following methods If you want a Page type and your item class name is Xxx
	// ====== template
	// type Xxx struct {
	// 		......
	// }
	// type XxxPage struct {
	// 		*internal.SdkPageable[*Xxx]
	// }
	// func NewXxxPageWithJsonString(str string) (*XxxPage, error) {
	// 	var o XxxPage
	// 	err := base.FromJsonString(str, &o)
	// 	return &o, err
	// }
}
