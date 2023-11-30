# golang2-fido-passkeys-api

> Passwordless FIDO Passkey API in golang, Passwordless future

[https://documenter.getpostman.com/view/12907432/2s9YeHYpym](Postman Collection)
[https://chat.openai.com/share/a5a947f0-9f3f-4045-aec8-967b8824c513](ChatGPT Prompts)

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
cd cli
go run . gen # generate RSA keys
go run . decrypt -c challenge-string # manually decrypt challenge string and solve manually too
go run . sign -m challenge-solution # sign the challenge solution
go run . register -n "User fullname" -e "user email" --server-url http://localhost:8080 # register user with previously generated rsa keys and verify challenge
go run . login --server-url http://localhost:8080 # login user using passkeyid
```

this creates `$HOME/.FIDO2` Folder with rsa keys and config.yml file
you can import or export keys in this folder

passkey.pem - private key
passkey.pub - public key
config.yml -  config file (not to be edited)
