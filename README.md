# scratch-test

A simple api endpoint with custom in-memory database implementation.

### Assumptions:
* A custom database is to be implemented and other existing database implementation
  can't be used.
  
* The data stored in the database is in the format as given in the `challenge` document

## Usage
1. Just running `go run main.go` from the `scratch-test` directory should get the
   server up and running.
    
2. A docker image can be built and run from the Dockerfile.
   a
   ```bash
   # Build the app image
   docker build -t scratch-app .
   
   # Run the app image in interactive mode to access logs
   docker run -it -p 8080:8080 scratch-app
   ```

3. A docker image is automatically pushed to `docker hub` on each push.
   The image can be pulled and used.
   ```bash
   docker run -p 8080:8080 lachsin/scratch-app:main
   ```

4. Only `/` endpoint is available, but different filters can be used.
    ```bash
    curl -X GET http://localhost:8080/?state=florida\&state=ak\&limit=4
    ```
   Returns four max entries that have state florida or alaska.

   Some of the available filters are:
   ```bash
   all=true     # Returns all the entry in the database
   state=ak&state=florida
   from=09:00&from=08:00
   to=15:00&to=20:00
   name=Dental Clinic&name=Vet Clinic
   limit=20     # Returns first 20 entries with given filters
   ```
   Paging is not supported.

5. Github Actions have been used as CI tool.

6. With the current amount of data in the database, the `jmeter` load
   testing showed:
    ```bash
    Request: http://localhost:8080?all=true
    Users: 10000
   
    Error: 0.00%,
    Throughput: 8887.3/s
    ```
   
