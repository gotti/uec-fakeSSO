module gotti/union

go 1.15

replace gotti/internal => ./internal

replace gotti/smtpMail => ./pkg/smtpMail

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	gotti/internal v0.0.0-00010101000000-000000000000 // indirect
	gotti/smtpMail v0.0.0-00010101000000-000000000000 // indirect
)
