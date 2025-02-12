# Environment variables

## Server

| Environment variable      | Description                              | Type   |
|---------------------------|------------------------------------------|--------|
| NAME                      | Name of the service                      | string |
| HOST                      | IP address or hostname of the service    | string |
| PORT                      | Port of the service                      | int    |
| PRODUCTION                | Production mode of the service           | bool   |
| GRACEFUL_SHUTDOWN_TIMEOUT | Graceful shutdown timeout of the service | time   |

## Logging

| Environment variable | Description                                                     | Type                        |
|----------------------|-----------------------------------------------------------------|-----------------------------|
| LOG_AS_JSON          | Logging this output as JSON. If deactivated, the output is text | bool                        |
| LOG_LEVEL            | Log level of the service                                        | DEBUG / INFO / WARN / ERROR |

## Recover

| Environment variable        | Description                             | Type  |
|-----------------------------|-----------------------------------------|-------|
| RECOVER_STACK_SIZE          | Set stack size in recovery              | int   |
| RECOVER_DISABLE_STACK_ALL   | Disable all stacks in recovery          | bool  |
| RECOVER_DISABLE_PRINT_STACK | Disable to print the stack in recovery  | bool  |

## Acl

| Environment variable | Description                                | Type   |
|----------------------|--------------------------------------------|--------|
| ACL_ENABLED          | Enables authentication                     | bool   |
| ACL_AUTH_USERNAME    | Username for authentication                | string |
| ACL_AUTH_PASSWORD    | Password for authentication                | string |
| ACL_AUTH_MODEL       | File to match the correct policy           | string |
| ACL_POLICY_MODEL     | File to define policy rules for api routes | string |

## RestClient

| Environment variable       | Description                                                    | Type   |
|----------------------------|----------------------------------------------------------------|--------|
| REST_CLIENT_BASE_URL       | Set global base url for all requests via the rest client       | string |
| REST_CLIENT_TIMEOUT        | Set global timeout for all requests via the rest client        | time   |
| REST_CLIENT_USERNAME       | Set global username for all requests via the rest client       | string |
| REST_CLIENT_PASSWORD       | Set global password for all requests via the rest client       | string |
| REST_CLIENT_TOKEN          | Set global token for all requests via the rest client          | string |
| REST_CLIENT_CONTENT_LENGTH | Set global content length for all requests via the rest client | bool   |

## Web Secure

| Environment variable          | Description                        | Type     |
|-------------------------------|------------------------------------|----------|
| SECURE_ENABLED                | Enables web secure                 | bool     |
| SECURE_HEADER_XSS             | Set xss header                     | string   |
| SECURE_HEADER_NO_SNIFF        | Set no sniff header                | string   |
| SECURE_HEADER_XFRAME          | Set xframe header                  | string   |
| SECURE_HEADER_MAX_AGE         | Set hsts max age header            | int      |
| SECURE_HEADER_CSP             | Set content security policy header | string   |
| SECURE_CORS_ALLOW_HEADERS     | Allow specified headers for cors   | []string |
| SECURE_CORS_ALLOW_METHODS     | Allow specified methods for cors   | []string |
| SECURE_CORS_ALLOW_ORIGINS     | Allow specified urls for cors      | []string |
| SECURE_CORS_ALLOW_CREDENTIALS | Allow to use credentials for cors  | bool     |
| SECURE_RATE_LIMIT             | Set rate limit                     | float64  |
| SECURE_RATE_BURST             | Set burst for rate limiter         | int      |
| SECURE_RATE_EXPIRES_IN        | Set expires in for rate limiter    | time     |
| SECURE_CSRF_TOKEN_LENGTH      | Set csrf token length              | uint8    |
| SECURE_CSRF_TOKEN_HEADER      | Set csrf token header              | string   |
| SECURE_CSRF_COOKIE_NAME       | Set csrf cookie name               | string   |
| SECURE_CSRF_COOKIE_MAX_AGE    | Set csrf cookie max age            | int      |
| SECURE_CSRF_COOKIE_SECURE     | Set csrf cookie secure             | bool     |

## MongoDB

| Environment variable          | Description                                             | Type          |
|-------------------------------|---------------------------------------------------------|---------------|
| MONGODB_HOST                  | IP address or hostname of the server                    | string        |
| MONGODB_PORT                  | Port of the server                                      | string        |
| MONGODB_DATABASE              | Database name of the server                             | string        |
| MONGODB_USERNAME              | Username of the server                                  | string        |
| MONGODB_PASSWORD              | Password of the server                                  | string        |
| MONGODB_COLLECTIONS | Collections to access                  | []string        |
| MONGODB_TLS_MODE              | TLS mode of the server                                  | string        |
| MONGODB_TLS_INSECURE          | TLS insecure of the server                              | string        |
| MONGODB_REPLICA_SET           | Replica set name -> used only on Cloud                  | string        |
| MONGODB_APP_NAME              | App name of the connection                              | string        |
| MONGODB_TIMEOUT               | Timeout of the connection                               | time.Duration |
| MONGODB_CONN_MAX_LIFETIME     | Max life time of the connection                         | time.Duration |
| MONGODB_MAX_IDLE_CONNECTIONS  | Max idle pool size of the connections                   | int           |
| MONGODB_MAX_OPEN_CONNECTIONS  | Max size of the open connections                        | int           |
| MONGODB_RETRY_WRITES          | Execute write operation once again after network errors | int           |
| MONGODB_DIRECT_CONNECTION     | Use direct connection via specified host                | int           |

## Postgres

| Environment variable          | Description                           | Type          |
|-------------------------------|---------------------------------------|---------------|
| POSTGRES_HOST                 | IP address or hostname of the server  | string        |
| POSTGRES_PORT                 | Port of the server                    | string        |
| POSTGRES_DATABASE             | Database name of the server           | string        |
| POSTGRES_USERNAME             | Username of the server                | string        |
| POSTGRES_PASSWORD             | Password of the server                | string        |
| POSTGRES_SSL_MODE             | SSL mode of the server                | string        |
| POSTGRES_SSL_CERT             | Certificate of the server             | string        |
| POSTGRES_MIGRATION_SOURCE     | Migration source of the server        | string        |
| POSTGRES_TIMEOUT              | Timeout of the connection             | time.Duration |
| POSTGRES_CONN_MAX_LIFETIME    | Max life time of the connection       | time.Duration |
| POSTGRES_MAX_IDLE_CONNECTIONS | Max idle pool size of the connections | int           |
| POSTGRES_MAX_OPEN_CONNECTIONS | Max size of the open connections      | int           |