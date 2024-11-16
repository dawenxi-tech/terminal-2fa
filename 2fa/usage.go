package main

const usage = `
usage:
	2fa - show totp code 
	2fa c[onfig] - config 2fa name and secret in command line
	2fa gui - show gui to config
`
const usageConfig = `
usage config:
	2fa c[onfig] [list|add|edit|delete|import|move] args...
usage config example:
	2fa config list
		list all 2fa
	2fa c add -name=foo -secret=bar
		add new 2fa when name is foo and secret is bar
	2fa config edit -id=1 -name=foo -secret=bar
		update the first 2fa, set name foo and set secret bar
	2fa c delete -id=1 -name=foo
		delete the first 2fa or delete which name is foo
	2fa config import -url=otpauth-migration://offline?data=CiYKCk1kb....
		import url which export by Google Authenticator
	2fa config move -id=1 -offset=2
		move the first 2fa to second position.
`
