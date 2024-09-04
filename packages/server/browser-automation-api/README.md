# Browser Automation Service (B.A.S)

## Table of Contents
- [About the Project](#about-the-project)
- [Features](#features)
- [Installation](#installation)
- [Developing](#developing)
- [API Documentation](#api-documentation)

## About the Project
CustomerOS Browser Automation Service (B.A.S) is an API service that allows you to automate browser tasks. It is built on top of the [Playwright](https://playwright.dev/) library. The service is designed to automate browser tasks programatically. It is useful for automating repetitive tasks such as form filling, data extraction, and more.

## Features
- Automate browser tasks
- Automate form filling
- Automate data extraction
- Automate data entry
- Automate browser navigation
- Automate browser interaction

## Installation
1. Clone the repo
```sh
git clone https://github.com/yourusername/project-name.git
```

2. Install NPM packages
```sh
npm install
```

3. Create a `.env` file in the root directory and add the environment variables found in the `.env.example` file
```sh
touch .env
```

## Developing
### Start the server in dev mode
```sh
npm run dev
```
### Generate a migration
```sh
npm run drizzle:generate
```
### Run the migration
```sh
npm run drizzle:migrate
```
### Push schema changes to the database without running a migration
This command is useful for prototyping database changes.
It will update the database schema without running a migration.
This command should not be used in production.
```sh
npm run drizzle:push
```

## API Documentation
The API documentation can be found [here](https://bas.customeros.ai/docs)
In development mode, the API documentation can be found at `http://localhost:3000/docs`
