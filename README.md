<h1 align="center">
  <br>
  <img src="./frontend/static/logo.svg" alt="The Final Testination" width="400" height="200">
</h1>

<h4 align="center">"The Final Testination" is an open-source gamified web application created by <a href="https://www.unige.it">Unige</a> students.</h4>


<p align="center">
  <a href="#key-key-features">Key Features</a>•
  <a href="#smile-set-up">Set up</a> •
  <a href="#runner-run-locally"> Run Locally </a>•
  <a href="#cop-tests-locally"> Test Locally </a>•
  <a href="#arrow_down-download">Download</a> •
  <a href="#lock-google-authentication">Google Authentication</a> •
  <a href="#email-emailware-da-cambiare">Emilware</a>•
  <a href="#mortar_board-license">License</a>
</p>

![Video non supportato](./frontend/static/The-Final-Testination-Video.gif)


## :key: Key Features

- **Login with Google**: Easily sign in using your Google account.
- **Email registration**: Register with a personal email address for access.
- **Level system**: 3 pre-existing levels with the option to add new ones.
- **Hints feature**: Get helpful hints to overcome challenges and progress.
- **User profile**: Track your completed levels and time spent.
- **Leaderboard**: Compare your performance with others on the leaderboard.

## :smile: Set Up

To clone and run this application, you'll need [Git](https://git-scm.com) and [sveltejs/kit](https://www.npmjs.com/package/@sveltejs/kit) which comes with [npm](http://npmjs.com). From your command line:

```bash
# Clone this repository
$ git clone https://github.com/Camo-02/The-Final-Testination

# Go into the repository
$ cd The-Final-Testination

# Install sveltejs/kit before running the front end
$ npm i @sveltejs/kit
```

## :runner: Run Locally
A `.env` file is required in the root directory of the project. A `.env.sample` is available in the same directory with some sensible default values.

```sh
#### Run this first in one terminal and wait for the DB to start...
docker compose -f development.yml up # DB

#### ...then run this in a second terminal...
cd backend && ../scripts/addenv go run main.go # Backend

#### ...and then this in the third and last terminal
cd frontend && ../scripts/addenv npm run dev # Frontend
```

### :hand: Stop the compose

This command will stop the database

```sh
docker compose -f development.yml down
```

## :cop: Test Locally

The following commands assume you are currently in the root of the project:

```sh
#### In one terminal
docker compose -f development.yml up

#### In another terminal
cd backend && ../scripts/addenv go test ./...

#### In a third terminal
cd frontend
npm i
../scripts/addenv npm run dev

#### In a fourth terminal
cd frontend
../scripts/addenv npm run build:test && ../scripts/addenv npm run host
# If you want to run all the tests from the cli run this in another terminal:
../scripts/addenv npm run cypress # add ` -- -b <browser-name>` (or path) 
# to run on a specific browser (e.g. `-- -b chrome`). 
# There's a specific shortcut for chrome: `npm run cypress:chrome`
# If you only want to run a specific type of tests, run:
../scripts/addenv npm run cypress:e2e # or cypress:component
# If you want to open cypress in a browser, then run
../scripts/addenv npm run cypress:open
```

> **Note**: you have to `npm i` and `npm run dev` in order to have the dependencies and the `.svelte-kit` directory before running cypress tests, otherwise they will not work.

### psql

To run queries as strings against the local db, run:

```sh
psql -U postgres -p 5432 -h 127.0.0.1 -d final_testination -c "[query]"
```
## :arrow_down: Download

## :lock: Google Authentication
To integrate Google authentication into your application, follow these steps:

### 1. Create OAuth Consent Screen
The OAuth consent screen is where users are prompted to grant access to their Google account information. This is where you define how your application will appear to users. You need to configure:

- **Application Name**
- **Authorized Domains**
- **Scopes** (permissions requested by the app)
- **Branding Information** (logo, application homepage)

### 2. Create Credentials
Once the OAuth consent screen is set up, you need to create credentials that will allow your app to authenticate users via Google:

- **Select OAuth 2.0 Client ID** as the credential type.
- Set the redirect URI for your application.
- Download the credentials file containing the client ID and client secret.

For further setup, visit the [Google Cloud Console](https://console.cloud.google.com).

### 3. Set Credentials in the frontend
Now you need to go to frontend/src/routes/login/+page.svelte and replace "Insert_Your_Client_ID" with the client id you obtained from the Google Console.

## :email: Emailware 

## :crown: Credits

## :mortar_board: License

MIT

---
