package repository

type IURLVerificationRepository interface {
	Challenge(s string) string
}

type URLVerificationRepository struct {
}

func NewURLVerificationRepository() *URLVerificationRepository{
	return &URLVerificationRepository{}
}
func (u *URLVerificationRepository) Challenge(s string) string {
	return s
}
