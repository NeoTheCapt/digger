package utils

type privilege struct {
	file   bool
	isroot bool
}

type credential struct {
	username string
	password string
	priv     privilege
}

func InitCred(username string, password string) credential {
	priv := privilege{false, false}
	return credential{
		username: username,
		password: password,
		priv:     priv,
	}
}
