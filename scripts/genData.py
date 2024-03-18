import faker
from random import choice, randint

filesLocation = "build/data/"

usersFile = filesLocation + "users.csv"
actorsFile = filesLocation + "actors.csv"
moviesFile = filesLocation + "movies.csv"
moviesActorsFile = filesLocation + "moviesActors.csv"

USERS_ROWS = 1000
ACTORS_ROWS = 1000
MOVIES_ROWS = 500
MOVIES_ACTORS_ROWS = 5000

myFaker = faker.Faker("en_US")


def generateUsers():
    role = ["user", "admin"]

    file = open(usersFile, "w", encoding="utf-8")

    for i in range(USERS_ROWS):
        line = "{};{};{};{}\n".format(i + 1, myFaker.unique.email(), myFaker.password(),
                                      choice(role))
        file.write(line)

    file.close()


def generateActors():
    gender = ["m", "f"]
    file = open(actorsFile, "w", encoding="utf-8")

    for i in range(ACTORS_ROWS):
        line = "{};{};{};{};{}\n".format(i + 1, myFaker.first_name(), myFaker.last_name(), choice(gender),
                                         myFaker.date_of_birth())

        file.write(line)

    file.close()


def generateMovies():
    file = open(moviesFile, "w", encoding="utf-8")

    for i in range(MOVIES_ROWS):
        line = "{};{};{};{};{}\n".format(i + 1, myFaker.text(max_nb_chars=20), myFaker.text(max_nb_chars=80), myFaker.date_this_decade(),
                                randint(0, 10))

        file.write(line)

    file.close()

def generateMoviesActors():
    file = open(moviesActorsFile, "w", encoding="utf-8")

    for i in range(MOVIES_ACTORS_ROWS):
        line = "{};{};{}\n".format(i + 1, randint(1, MOVIES_ROWS), randint(1, ACTORS_ROWS))

        file.write(line)

    file.close()


generateUsers()
generateActors()
generateMovies()
generateMoviesActors()
