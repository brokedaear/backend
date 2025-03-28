package dal

type S3 struct{}

func NewS3Storage() *S3 {
	return &S3{}
}
