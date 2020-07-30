# EVM Backend Assessment

This repository contains all code that is submitted as part of Evermos Backend
Assessment.

## 01. Machine Gun

This is supposed to be the solution to Question 1 of Evermos Backend Engineer
Assessment v1.1. **I later abandoned this in favor of v1.2**.

## 02. Kitara Store API

> Running Locally Requires:
> * Go 1.14+
> * [direnv](https://direnv.net/)

This is the solution to **Question 2 of Evermos Backend Engineer Assessment
v1.2**, which is identical to **Question 2 of Evermos Backend Engineer
Assessment v1.1**. This API is not available publicly and is deliberately
minimal to emphasize on implementing concurrency handling as indicated in the
problem statement.

To run locally:

1. Clone this repository.
2. Go to `02-kitara-store` folder.
3. Create `.envrc` file based on the sample file, then make sure it is allowed
   (run `direnv allow`).
4. Configure your database and run the migrations in `migrations` folder.
3. `go run main.go`

The source code includes functional tests that are run concurrently to showcase
the concurrency handling of the API.

To test the concurrency feature:

1. Have the API running as explained above.
2. Open up a new terminal, go to `tests/functional` folder, then run `go test`.
3. The test should show no errors.
4. To repeat the test, first clear the database and re-run the migrations.

## 03. Key Puzzle

This is the solution to **Question 3 of Evermos Backend Engineer Assessment
v1.1**. To run locally, simply go into `03-key-puzzle` directory and invoke
`go run main.go`

## 04. Tennis Player API

> Running Locally Requires:
> * Go 1.14+
> * [direnv](https://direnv.net/)

> Public addresses:
> * API: http://tennis.evm.radityakertiyasa.com:8080/
> * Swagger Docs: http://tennis.evm.radityakertiyasa.com:8080/docs/index.html

This is the solution to **Question 1 of Evermos Backend Engineer Assessment
v1.2**. This API is available publicly and testing can be done using the
included Swagger Docs UI.

## Assumptions

To produce a reasonable API and logic within it, I made the following
assumptions based on the problem statement:

* Rahman represents a Player entity.
* A Player can have many Containers.
* A Container will have a set capacity, and once it is full it cannot receive
  any more balls. This represents the "verified" state of the Container.
* The Player can put a ball randomly into any of his Containers.
* Once a Container is full (verified), the Player is then ready to play and he
  cannot put any more balls into any of his Containers.

## Testing the Functionality

To test the API's functionality, I have included Swagger Docs UI that is
available [here](http://tennis.evm.radityakertiyasa.com:8080/docs/index.html).
The steps to test are as follows:

* Create a new Player.
* Create new Containers that belong to that new Player. Create as many as you'd
  like.
* Invoke the `addBall` endpoint and observe the result. A single ball will
  added randomly into any one of the available Containers attached to the
  Player. Do this until one of the Containers is full
* Once one of the Containers is full, you cannot invoke `addBall` on a Player
  because he is now ready to play.
