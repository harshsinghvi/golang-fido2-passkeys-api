# golang2-fido-passkeys-api

> Passwordless FIDO Passkey API in golang, Passwordless future

- <https://chat.openai.com/share/a5a947f0-9f3f-4045-aec8-967b8824c513> - ChatGPT Prompts
- <https://documenter.getpostman.com/view/12907432/2s9YeHYpym> - Postman Collection

## Documentation to be updated soon

- WIP: Documentation

## Usecases

- CLI Apps suthentication like ssh
- Mobile based passkeys (Passwordless authentication)

## CLI Usage

```bash
go install github.com/harshsinghvi/golang-fido2-passkeys-api/cli # install locally after cloning
go install github.com/harshsinghvi/golang-fido2-passkeys-api/cli@latest # install directly

cli decrypt -c challenge-string # manually decrypt challenge string and solve manually too
cli sign -m challenge-solution # sign the challenge solution

cli gen # generate RSA keys
cli register -n "User fullname" -e "user email" --server-url http://localhost:8080 # register user with previously generated rsa keys and verify challenge
cli register-new-key -e email -d description --server-url http://localhost:8080 # add key to user account
cli login --server-url http://localhost:8080 # login user using stored keys
cli get-me # Business logic
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

- test ReVerify user and passkey endpoint
- check public key encoding

- orgs

- error handeling
- rate limit
- database events
- user preferneces
- check BillingDisable

- Workflow for releasing binaries
- token roles
- clean code in `cli` and `crypto` library

- make new repos for cli and `autoroutes` routes
UI
- admin portal
- org login
- WIP: Documentation

### TEST DB

#### -- make changes in this

```sql
CREATE DATABASE test_db_savepoint; 
```

#### create test db from savepoint

```sql
SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity 
WHERE pg_stat_activity.datname in ('test_db_savepoint' ,'test_db') AND pid <> pg_backend_pid();

CREATE DATABASE test_db WITH TEMPLATE test_db_savepoint OWNER postgres;
```

#### reset to savepoint

```sql
SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity 
WHERE pg_stat_activity.datname in ('test_db_savepoint' ,'test_db') AND pid <> pg_backend_pid();
DROP DATABASE test_db;
CREATE DATABASE test_db WITH TEMPLATE test_db_savepoint OWNER postgres;
```

#### hard delete

```sql
DELETE FROM access_logs WHERE deleted_at IS NOT NULL;
DELETE FROM access_tokens WHERE deleted_at IS NOT NULL;
DELETE FROM challenges WHERE deleted_at IS NOT NULL;
DELETE FROM events WHERE deleted_at IS NOT NULL;
DELETE FROM passkeys WHERE deleted_at IS NOT NULL;
DELETE FROM users WHERE deleted_at IS NOT NULL;
DELETE FROM verifications WHERE deleted_at IS NOT NULL;
```
