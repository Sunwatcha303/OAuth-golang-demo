# OAuth-golang-demo
demo Open Authentication with golang

## Create .env 
put **.env** on the root path

``` 
FIBER_HOST=0.0.0.0
FIBER_PORT=5001

DB_HOST=0.0.0.0
DB_PORT=1122
DB_DATABASE=your_database_name
DB_USERNAME=your_database_username
DB_PASSWORD=our_database_password
DB_PROTOCOL=tcp
DB_SSL_MODE=disable

CLIENT_ID=your_client_id
CLIENT_SECRET=your_client_secret
REDIRECT_URI=http://localhost:5173/callback

REDIS_HOST=0.0.0.0
REDIS_PORT=6379
REDIS_DATABASE=0

JWR_SECRET_KEY=your_jwt_secret_key

```
## Start the server
use Makefile to run app
``` 
make

```
