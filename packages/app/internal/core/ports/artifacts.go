package ports

type ArtifactStorage interface {
	Get()
	Add()
	Tokenize()
}
