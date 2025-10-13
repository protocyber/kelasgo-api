# kelasgo-api

## Description

This is a template for your project's readme file. Change it as necessary to
provide basic information about this project. For more detailed documentation,
use Confluence or the project's Wiki pages.

### Setup

1. Clone repo
1. Copy and paste `.env.example` to `.env`. Set your database and the other configs
1. Make sure you install [go-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)
1. Run the migration : `make migrate_up`
1. Run the service : `make dev`

### Database migration

As you know above, we are using [go-migrate](https://github.com/golang-migrate/migrate) for the migration tool. We have simplify the frequently executed command into the Makefile.

**Note : make sure you set these env properly :** 
```
DB.MYSQL.WRITE.USER
DB.MYSQL.WRITE.PASSWORD
DB.MYSQL.WRITE.HOST
DB.MYSQL.WRITE.PORT
DB.MYSQL.WRITE.NAME
```
We're use those environment variable to open the database connection to perform the migration.

| Command | Description |
| --- | --- |
| `make migrate_create` | Create migration file |
| `make migrate_up` | Apply all or N up migrations. If you want to specify the step, please use : `make migrate_up MIGRATION_STEP=<number>` |\
| `make migrate_down` | Apply all or N down migrations. If you want to specify the step, please use : `make migrate_down MIGRATION_STEP=<number>`
| `make migrate_force` | Set version V but don't run migration (fix the dirty state) |
| `make migrate_version` | Print current migration version |
| `make migrate_drop` | Drop everything inside database |

More info : https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

## Management

### Communication Channels

Please use the following channels for communications about this project.

* Discord Channel(s):
  * [#CHANNEL1]
  * [#CHANNEL2]
* JIRA Board(s):
  * [JIRA-BOARD-1]
  * [JIRA-BOARD-2]

### Product Managers

Please contact the persons below for product and management inquries.

* [PM1]
* [PM2]

### Engineering Managers

Please contact the persons below for engineering management inquries.

* [EM1]
* [EM2]

## Maintainers

Please contact the persons below for techincal inquries.

### Tech Lead

[TL1]

### Frontend Engineers

* [FE1]
* [FE2]

### Backend Engineers

* [BE1]
* [BE2]

### Quality Assurance Engineers

* [QAE1]
* [QAE2]

## Technical Information

For technical information about Go-based Backend Services, please visit the
following links:
