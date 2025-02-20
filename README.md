# JWT Authentication

JWT Aunthentication and Authorization using golang

## Technologies Used

1. Go (Golang) - Backend language
2. Fiber - Web framework for Go
3. MongoDB - NoSQL database for storing user data
4. JWT (JSON Web Tokens) - Token-based authentication
5. Docker - Containerization for running MongoDB and API

## Running the Project

Essentials:

1. Make sure you have docker installed on your system
2. Clone the repository and open the project

Run the below command

```
docker-compose up --build
```

## API Endpoints and Usage

1. Sign Up:
   Registers a new user with an email and password.

```curl
curl -X POST http://localhost:3000/auth/signup \
     -H "Content-Type: application/json" \
     -d '{
           "email": "user@example.com",
           "password": "securepassword"
         }'
```

2. Sign In:
   Authenticates a user and returns a JWT token.

```curl
curl -X POST http://localhost:3000/auth/signin \
     -H "Content-Type: application/json" \
     -d '{
           "email": "user@example.com",
           "password": "securepassword"
         }'
```

3. Refresh Token:
   Generates a new token using the old one.

```curl
curl -X GET http://localhost:3000/auth/refresh \
     -H "Authorization: Bearer ${Token}"
```

4. Revoke Token:
   Invalidates the provided JWT token.

```curl
curl -X POST http://localhost:3000/auth/revoke \
     -H "Authorization: Bearer ${Token}"
```

5. Get User Details:
   Fetches the details of the authenticated user.

```curl
curl -X GET http://localhost:3000/user/details \
     -H "Content-Type: application/json"  \
     -H "Authorization: Bearer ${Token}"
```

Note:
In the above-mentioned curl requests replace your obtained token with `${Token}`

## License

[MIT](https://choosealicense.com/licenses/mit/)
