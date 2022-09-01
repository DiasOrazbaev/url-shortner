package shortner

type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Save(redirect *Redirect) error
}
