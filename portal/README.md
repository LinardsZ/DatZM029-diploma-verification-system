# `${REPO_NAME_TITLE}`

<!-- TOC -->

- [${REPO_NAME_TITLE}](#repo_name_title)
  - [Development](#development)
  - [Requirements](#requirements)
  - [Customize configuration](#customize-configuration)
  - [:warning: IMPORTANT NOTES :warning:](#warning-important-notes-warning)
  - [Testing docker image locally](#testing-docker-image-locally)

<!-- /TOC -->

## Development

1. Build dev server:

    ```bash
    pnpm i
    ```

2. Run dev server (_also possible with vs code debug functionality (F5)_) :

    ```bash
    pnpm dev
    ```

## Requirements

- [Node.js](https://nodejs.org/en/) (at least v20.14.0)
- [PNPM](https://pnpm.io/) (at least v10.7.0 )
- [ni](https://github.com/antfu/ni) (optional - in order to use `ni` command instead of `pnpm`)

## Customize configuration

Environment variables are loaded from `.env` files in the root directory. See [Vite Environment Variables and Modes](https://vitejs.dev/guide/env-and-mode.html) for more information.

| Variable | Description | Default on serve (locally) | Default on build |
| --- | --- | --- | --- |
| `SERVICE_URL` | URL of the API service that your application will use | '' | /api |
| `AUTH_URL` | URL of the authentication service that your application will use | '' | /idauth |
| `ENVIRONMENT` | Environment name | development | production |
| `PUBLIC_URL` | Public URL of the application. Only port 44342 is supported if you want to test against lx-demo api. | <https://localhost:44342/> |  |
| `BASE_PATH` | Base path of the application | / | / |
| `CLIENT_ID` | VPM client ID | ${REPO_NAME_LOWER} | ${REPO_NAME_LOWER} |
| `USE_MOCK_MIDDLEWARE` | Use mock middleware [mocks/middleware.mock.js](./mocks/middleware.mock.js) if value `true` | false | always false (ignores env value) |

example .env file:

```sh
SERVICE_URL=https://localhost:44342/
AUTH_URL=https://localhost:44342/auth
ENVIRONMENT=development
PUBLIC_URL=https://localhost:44342/
USE_MOCK_MIDDLEWARE=false
```

## :warning: IMPORTANT NOTES :warning:

It's **very important** that the webapp (this portal) can be built and run locally as described below, with no extra steps:

- Clone this repo;
- Add `.env` file to my local project (see above);
- Run `pnpm i`;
- Press **F5** ("Run and Debug" in VSCode);
- Webapp starts in my local browser and can call published Dev API (without any CORS, HTTPS, Authentication redirect and/or other problems);

Other ways to run this webapp (e.g., connecting to locally run API, or connecting to test env, etc) are permitted, of course, but should be considered **additional** methods.

