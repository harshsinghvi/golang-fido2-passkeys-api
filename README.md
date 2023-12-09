# golang2-fido-passkeys-api

> Passwordless FIDO Passkey API in golang, Passwordless future

- <https://chat.openai.com/share/a5a947f0-9f3f-4045-aec8-967b8824c513> - ChatGPT Prompts
- <https://documenter.getpostman.com/view/12907432/2s9YeHYpym> - Postman Collection

- API
- Client CLI

## Documentation to be updated soon

- WIP: Documentation
- WIP: Api more api endpints
- WIP: Updated CLI
- WPI: Business Logic for demonstration

## Usecases

- CLI Apps suthentication like ssh
- Mobile based passkeys (Passwordless authentication)

## CLI Usage

```bash
go install github.com/harshsinghvi/golang-fido2-passkeys-api/cli # install locally after cloning
go install github.com/harshsinghvi/golang-fido2-passkeys-api/cli@latest # install directly

cli gen # generate RSA keys
cli decrypt -c challenge-string # manually decrypt challenge string and solve manually too
cli sign -m challenge-solution # sign the challenge solution
cli register -n "User fullname" -e "user email" --server-url http://localhost:8080 # register user with previously generated rsa keys and verify challenge
cli login --server-url http://localhost:8080 # login user using stored keys
cli get-me # Business logic
cli add-key -e email -d description --server-url http://localhost:8080 # add key to user account
```

this creates `$HOME/.FIDO2` Folder with rsa keys and config.yml file
you can import or export keys in this folder

- passkey.pem - private key
- passkey.pub - public key
- config.yml -  config file (not to be edited)

## Build Multi Arch Binary for CLI and Server

- `./go-executable-build.bash github.com/harshsinghvi/golang-fido2-passkeys-api`
- reference script: <https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04>
- go supported arch and os <https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63>

## TODO

- orgs
- nil check arrays in whole code
- send email

- Workflow for releasing binaries
- public key validation
- token roles
- ~~delete or disabled endpoint~~
- ~~verification expiry and updated verification process~~
- ~~user roles~~
- ~~user profile~~
- ~~email validation~~
- ~~demo business logic (todo app)~~
- ~~demo business logic cli~~
- ~~email verification after user creation~~
- ~~new passkey registration (followed by authorisation using email verification)~~
- ~~audit logs~~
- admin portal
- org login
