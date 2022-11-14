package categories

type Repository interface {
	FindById(id uint) (category Category, err error)
}
