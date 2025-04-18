# github.com/dhax/go-base MOTO Server

MOTO REST API for RFID-based system.

## Routes

<details>
<summary>`/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Timeout](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/**
	- _GET_
		- [SPAHandler](/api/api.go#L101)

</details>
<details>
<summary>`/admin/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Timeout](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/admin/**
	- [RequiresRole](/auth/authorizer.go#L11)
	- **/**
		- _GET_
			- Admin dashboard root

</details>
<details>
<summary>`/admin/*/accounts/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Timeout](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/admin/**
	- [RequiresRole](/auth/authorizer.go#L11)
	- **/accounts/**
		- **/**
			- _GET_
				- List all accounts
			- _POST_
				- Create a new account

</details>
<details>
<summary>`/admin/*/accounts/*/{accountID}/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Timeout](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/admin/**
	- [RequiresRole](/auth/authorizer.go#L11)
	- **/accounts/**
		- **/{accountID}/**
			- **/**
				- _PUT_
					- Update account
				- _DELETE_
					- Delete account
				- _GET_
					- Get account details

</details>
<details>
<summary>`/api/*/account/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Timeout](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/api/**
	- **/account/**
		- **/**
			- _PUT_
				- Update own account
			- _DELETE_
				- Delete own account
			- _GET_
				- Get own account details

</details>
<details>
<summary>`/auth/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Timeout](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/auth/**
	- **/login**
		- _POST_
			- Request a login token (passwordless auth)
	- **/token**
		- _POST_
			- Exchange login token for JWT access token
	- **/refresh**
		- _POST_
			- Refresh JWT token
	- **/logout**
		- _POST_
			- Invalidate JWT token

</details>
<details>
<summary>`/rfid/*`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Timeout](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/rfid/**
	- **/tag**
		- _POST_
			- Submit RFID tag read from Python daemon
	- **/tags**
		- _GET_
			- Get all stored RFID tags
	- **/app/sync**
		- _POST_
			- Sync tags from Tauri app
	- **/app/status**
		- _GET_
			- Get server status for Tauri app

</details>
<details>
<summary>`/healthz`</summary>

- [Recoverer](/vendor/github.com/go-chi/chi/middleware/recoverer.go#L18)
- [RequestID](/vendor/github.com/go-chi/chi/middleware/request_id.go#L63)
- [Timeout](/vendor/github.com/go-chi/chi/middleware/timeout.go#L33)
- [Logger](/vendor/github.com/go-chi/chi/middleware/logger.go#L36)
- [SetContentType](/vendor/github.com/go-chi/render/content_type.go#L49)
- **/healthz**
	- _GET_
		- Health check endpoint
</details>
