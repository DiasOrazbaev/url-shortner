package shortner

type RedirectService interface {
	Find(code string) (*Redirect, error)
	Save(redirect *Redirect) error
}
