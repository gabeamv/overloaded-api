# overloaded-api

Overloaded API is a fitness web api where users can track their progress and [progressively overload](https://en.wikipedia.org/wiki/Progressive_overload) for each
of their workouts.

## Features

* Authentication
    - Users can register an account and login. Refresh and access tokens are used for authenticated users to make requests to server resources.
* Exercises
    - Users can create and delete their own custom exercises. Owners and devs for the API can create default exercises used by all users.
* Workouts
    - Users can upload their workouts, keeping track of when it was done, how many PRs they accomplished, as well as the total volume that has been done.
* Workout Sets
    - Users upload workout sets in batches, which are tied to a specific workout that they recorded. Weight, reps, and time are recorded to describe the set. Users must also record a set with its progression rule (as seen below) which is a character that has its own unique meaning that signifies progression.
* Progression Rule
    - Progressive overload is the main driver of muscle growth. In order to build muscle, you need to try to push yourself harder a little after each workout, whether its adding more weight, doing more reps, or improving form. Currently there are 3 main rules (you can add on to it later if you want to):
        * p (progress)
        * s (stay)
        * t (tag on form)
    Choose a weight and a rep range for an exercise. For the set you are doing for that exercise, aim for the top end of the rep range. If you can achieve that in your set, mark that set as 'p', which means to progress in weight next workout in the same rep range. If you are unable to achieve the top end of the rep range, mark the set as 's' which means to stay at the current weight you are at and try again next workout. If you notice that you may be cheating on some reps or not accomplishing the set with decent enough form throughout the rep range, mark the set as 't'. This means you have tag on your form, so don't progress in weight. You can accomplish well past the rep range, but if you feel as if your form wasn't as intact, stay at the same weight and try again next workout. Ideally, the sets you want to track are the sets you push to failure, true working sets.

## Requirements

* Go 1.21+ installed
* A terminal to run commands

## Installation

1. Clone the repository.
2. Install and setup PostgreSQL. 
    * On mac, run **brew install postgres@15** or another suitable version.
    * On Linux / WSL (Debian), run **sudo apt update**, then **sudo apt install postgresql postgresql-contrib**
3. Ensure you have the postgres by running **psql --version**
    * If you are on Linux, update postgres password: **sudo passwd postgres**
    * Don't forget the password.
4. Enter the psql shell: **psql postgres**
5. Create a new database called '**overloaded**': **CREATE DATABASE overloaded;**
6. Set the user password (Linux only): **ALTER USER postgres PASSWORD 'postgres';**
7. In the root of the cloned project, we are going to add a '.env' file with some values: 
    1. DB_URL = "postgres://{username}:{password}@localhost:5432/overloaded?sslmode=disable"
        * Replace {username} and {password} with your own configuration
    2. PLATFORM = "dev"
    3. SECRET = "{generated secret}"
        * This secret will be used to sign and verify JWT tokens. To generate a secret, go into your terminal and run **openssl rand -base64 64**
        * Assign the generated value to 'SECRET'.

## Run

1. In terminal, run **go build && ./overloaded-api**.
2. To make requests to this api, you can utilize [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client), [curl](https://curl.se/), your own built client, or any other client. See the paths below for making requests.

## API Paths

[Docs](./documentation/api_documentation.md)