# Reproducible integration test environment

## Testing process

1. Start up Docker Compose stack
2. DB container will initialize with a dump we previously created (see below)
3. Healthchecks container will use that DB data
4. We can now run tests with known API endpoints, keys etc.

### Creating the DB dump

#### Using VSCode

1. Comment out the volume binding
2. Run the `Start test environment` task
3. Configure your Healthchecks instance as required, saving UUIDs/slugs and ping keys to `.env`
4. Run the `Create DB dump` task
5. Run the `Stop test environment` task

#### Manually

1. Boot up the stack from scratch by **commenting out the volume binding to `/docker-entrypoint-initdb.d`**
2. Configure your Healthchecks instance as required, saving UUIDs/slugs and ping keys to `.env`
3. Run `pg_dump -U <username> -d <database> > /healthchecks_dump.sql`
4.  Export the dump by running `docker cp <containerId>:/healthchecks_dump.sql ./healthchecks_dump.sql`
