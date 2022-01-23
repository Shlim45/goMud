# goMud

A MUD Engine written in goLang.

Uses MySQL 8 for database.
* Service must be running on _localhost_:_3306_ at startup.
* A `.env` file in the in the project root directory must specify:
  * DB_HOST - domain/ip where mysql is running
  * DB_PORT
  * DB_NAME - name of the database
  * DB_USER - username
  * DB_PASSWORD
* Required tables are automatically created if they do not exist in the specified database.