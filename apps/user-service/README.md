# User Service Guide

NOTE: The User Service is adapted from [here](https://github.com/CS3219-AY2425S1/PeerPrep-UserService/tree/main).

The User Service is a microservice responsible for handling user-related operations, such as registration, authentication, profile management, and authorization. It is part of a larger system architecture and integrates with other microservices within the project. Below outlines the requirements, architecture, and design for the User Service, and provides a thorough overview of its functionality, structure, and technology stack.

The User Service is designed using a microservices architecture, integrating with other
services through RESTful APIs. Below is a description of the architecture components
and their interactions.

1. User Service: This microservice handles all user-related operations, including user creation, authentication, data retrieval, updates, and role management. It exposes REST endpoints that allow clients to interact with the database.
2. MongoDB: MongoDB is used as the database to store user information. It is hosted on MongoDB Cloud.
3. Authentication: JWT tokens are generated upon login to authenticate users. This service ensures secure access to data based on user roles and verifies token validity for routing.

### Technology Stack

- Node.js: The primary runtime environment for the User Service due to its scalability, ease of use for API development, and compatibility with MongoDB.
- MongoDB Cloud: Used for its scalability and ease of management. MongoDB’s document-based structure is well-suited for handling user data.
- JWT (JSON Web Tokens): Used for secure and stateless authentication.

### Authentication and Authorization

The User Service uses JWT for secure token-based authentication. Upon login, a JWT is issued.

This token is used in the Authorization header to grant or restrict access to specific routes based on the user’s role:

- Admin Users: Can access all routes, including those that involve managing other users’ data, or editing question details.
- Non-admin Users: Restricted to actions related to their own data only.

---

NOTE: This following parts in the README is adapted from [here](https://github.com/CS3219-AY2425S1/PeerPrep-UserService/tree/main).

## Setting-up

> :notebook: If you are familiar to MongoDB and wish to use a local instance, please feel free to do so. This guide utilizes MongoDB Cloud Services.

1. Set up a MongoDB Shared Cluster by following the steps in this [Guide](./MongoDBSetup.md).

2. After setting up, go to the Database Deployment Page. You would see a list of the Databases you have set up. Select `Connect` on the cluster you just created earlier on for User Service.

   ![alt text](./GuideAssets/ConnectCluster.png)

3. Select the `Drivers` option, as we have to link to a Node.js App (User Service).

   ![alt text](./GuideAssets/DriverSelection.png)

4. Select `Node.js` in the `Driver` pull-down menu, and copy the connection string.

   Notice, you may see `<password>` in this connection string. We will be replacing this with the admin account password that we created earlier on when setting up the Shared Cluster.

   ![alt text](./GuideAssets/ConnectionString.png)

5. In the `user-service` directory, create a copy of the `.env.sample` file and name it `.env`.

6. Update the `DB_CLOUD_URI` of the `.env` file, and paste the string we copied earlier in step 4. Also remember to replace the `<password>` placeholder with the actual password.

## Running User Service

1. Open Command Line/Terminal and navigate into the `user-service` directory.

2. Run the command: `npm install`. This will install all the necessary dependencies.

3. Run the command `npm start` to start the User Service in production mode, or use `npm run dev` for development mode, which includes features like automatic server restart when you make code changes.

4. Using applications like Postman, you can interact with the User Service on port 3001. If you wish to change this, please update the `.env` file.

## User Service API Guide

### Create User

- This endpoint allows adding a new user to the database (i.e., user registration).

- HTTP Method: `POST`

- Endpoint: http://localhost:3001/users

- Body

  - Required: `username` (string), `email` (string), `password` (string)

    ```json
    {
      "username": "SampleUserName",
      "email": "sample@gmail.com",
      "password": "SecurePassword"
    }
    ```

- Responses:

  | Response Code               | Explanation                                           |
  | --------------------------- | ----------------------------------------------------- |
  | 201 (Created)               | User created successfully, created user data returned |
  | 400 (Bad Request)           | Missing fields                                        |
  | 409 (Conflict)              | Duplicate username or email encountered               |
  | 500 (Internal Server Error) | Database or server error                              |

### Get User

- This endpoint allows retrieval of a single user's data from the database using the user's ID.

  > :bulb: The user ID refers to the MongoDB Object ID, a unique identifier automatically generated by MongoDB for each document in a collection.

- HTTP Method: `GET`

- Endpoint: http://localhost:3001/users/{userId}

- Parameters

  - Required: `userId` path parameter
  - Example: `http://localhost:3001/users/60c72b2f9b1d4c3a2e5f8b4c`

- <a name="auth-header">Headers</a>

  - Required: `Authorization: Bearer <JWT_ACCESS_TOKEN>`

  - Explanation: This endpoint requires the client to include a JWT (JSON Web Token) in the HTTP request header for authentication and authorization. This token is generated during the authentication process (i.e., login) and contains information about the user's identity. The server verifies this token to ensure that the client is authorized to access the data.

  - Auth Rules:

    - Admin users: Can retrieve any user's data. The server verifies the user associated with the JWT token is an admin user and allows access to the requested user's data.
    - Non-admin users: Can only retrieve their own data. The server checks if the user ID in the request URL matches the ID of the user associated with the JWT token. If it matches, the server returns the user's own data.

- Responses:

  | Response Code               | Explanation                                              |
  | --------------------------- | -------------------------------------------------------- |
  | 200 (OK)                    | Success, user data returned                              |
  | 401 (Unauthorized)          | Access denied due to missing/invalid/expired JWT         |
  | 403 (Forbidden)             | Access denied for non-admin users accessing others' data |
  | 404 (Not Found)             | User with the specified ID not found                     |
  | 500 (Internal Server Error) | Database or server error                                 |

### Get All Users

- This endpoint allows retrieval of all users' data from the database.
- HTTP Method: `GET`
- Endpoint: http://localhost:3001/users
- Headers

  - Required: `Authorization: Bearer <JWT_ACCESS_TOKEN>`
  - Auth Rules:

    - Admin users: Can retrieve all users' data. The server verifies the user associated with the JWT token is an admin user and allows access to all users' data.
    - Non-admin users: Not allowed access.

- Responses:

  | Response Code               | Explanation                                      |
  | --------------------------- | ------------------------------------------------ |
  | 200 (OK)                    | Success, all user data returned                  |
  | 401 (Unauthorized)          | Access denied due to missing/invalid/expired JWT |
  | 403 (Forbidden)             | Access denied for non-admin users                |
  | 500 (Internal Server Error) | Database or server error                         |

### Update User

- This endpoint allows updating a user and their related data in the database using the user's ID.

- HTTP Method: `PATCH`

- Endpoint: http://localhost:3001/users/{userId}

- Parameters

  - Required: `userId` path parameter

- Body

  - At least one of the following fields is required: `username` (string), `email` (string), `password` (string)

    ```json
    {
      "username": "SampleUserName",
      "email": "sample@gmail.com",
      "password": "SecurePassword"
    }
    ```

- Headers

  - Required: `Authorization: Bearer <JWT_ACCESS_TOKEN>`
  - Auth Rules:

    - Admin users: Can update any user's data. The server verifies the user associated with the JWT token is an admin user and allows the update of requested user's data.
    - Non-admin users: Can only update their own data. The server checks if the user ID in the request URL matches the ID of the user associated with the JWT token. If it matches, the server updates the user's own data.

- Responses:

  | Response Code               | Explanation                                             |
  | --------------------------- | ------------------------------------------------------- |
  | 200 (OK)                    | User updated successfully, updated user data returned   |
  | 400 (Bad Request)           | Missing fields                                          |
  | 401 (Unauthorized)          | Access denied due to missing/invalid/expired JWT        |
  | 403 (Forbidden)             | Access denied for non-admin users updating others' data |
  | 404 (Not Found)             | User with the specified ID not found                    |
  | 409 (Conflict)              | Duplicate username or email encountered                 |
  | 500 (Internal Server Error) | Database or server error                                |

### Update User Privilege

- This endpoint allows updating a user’s privilege, i.e., promoting or demoting them from admin status.

- HTTP Method: `PATCH`

- Endpoint: http://localhost:3001/users/{userId}

- Parameters

  - Required: `userId` path parameter

- Body

  - Required: `isAdmin` (boolean)

    ```json
    {
      "isAdmin": true
    }
    ```

- Headers

  - Required: `Authorization: Bearer <JWT_ACCESS_TOKEN>`
  - Auth Rules:

    - Admin users: Can update any user's privilege. The server verifies the user associated with the JWT token is an admin user and allows the privilege update.
    - Non-admin users: Not allowed access.

> :bulb: You may need to manually assign admin status to the first user by directly editing the database document before using this endpoint.

- Responses:

  | Response Code               | Explanation                                                     |
  | --------------------------- | --------------------------------------------------------------- |
  | 200 (OK)                    | User privilege updated successfully, updated user data returned |
  | 400 (Bad Request)           | Missing fields                                                  |
  | 401 (Unauthorized)          | Access denied due to missing/invalid/expired JWT                |
  | 403 (Forbidden)             | Access denied for non-admin users                               |
  | 404 (Not Found)             | User with the specified ID not found                            |
  | 500 (Internal Server Error) | Database or server error                                        |

### Delete User

- This endpoint allows deletion of a user and their related data from the database using the user's ID.
- HTTP Method: `DELETE`
- Endpoint: http://localhost:3001/users/{userId}
- Parameters

  - Required: `userId` path parameter

- Headers

  - Required: `Authorization: Bearer <JWT_ACCESS_TOKEN>`

  - Auth Rules:

    - Admin users: Can delete any user's data. The server verifies the user associated with the JWT token is an admin user and allows the deletion of requested user's data.

    - Non-admin users: Can only delete their own data. The server checks if the user ID in the request URL matches the ID of the user associated with the JWT token. If it matches, the server deletes the user's own data.

- Responses:

  | Response Code               | Explanation                                             |
  | --------------------------- | ------------------------------------------------------- |
  | 200 (OK)                    | User deleted successfully                               |
  | 401 (Unauthorized)          | Access denied due to missing/invalid/expired JWT        |
  | 403 (Forbidden)             | Access denied for non-admin users deleting others' data |
  | 404 (Not Found)             | User with the specified ID not found                    |
  | 500 (Internal Server Error) | Database or server error                                |

### Login

- This endpoint allows a user to authenticate with an email and password and returns a JWT access token. The token is valid for 1 day and can be used subsequently to access protected resources. For example usage, refer to the [Authorization header section in the Get User endpoint](#auth-header).
- HTTP Method: `POST`
- Endpoint: http://localhost:3001/auth/login
- Body

  - Required: `email` (string), `password` (string)

    ```json
    {
      "email": "sample@gmail.com",
      "password": "SecurePassword"
    }
    ```

- Responses:

  | Response Code               | Explanation                                        |
  | --------------------------- | -------------------------------------------------- |
  | 200 (OK)                    | Login successful, JWT token and user data returned |
  | 400 (Bad Request)           | Missing fields                                     |
  | 401 (Unauthorized)          | Incorrect email or password                        |
  | 500 (Internal Server Error) | Database or server error                           |

### Verify Token

- This endpoint allows one to verify a JWT access token to authenticate and retrieve the user's data associated with the token.
- HTTP Method: `GET`
- Endpoint: http://localhost:3001/auth/verify-token
- Headers

  - Required: `Authorization: Bearer <JWT_ACCESS_TOKEN>`

- Responses:

  | Response Code               | Explanation                                        |
  | --------------------------- | -------------------------------------------------- |
  | 200 (OK)                    | Token verified, authenticated user's data returned |
  | 401 (Unauthorized)          | Missing/invalid/expired JWT                        |
  | 500 (Internal Server Error) | Database or server error                           |

### Build Docker

```bash
# Navigate to the user-service app directory
cd apps/user-service

# Build dockerfile (Ensure that your docker daemon is running beforehand)
docker build -t user-service -f Dockerfile .
```

### Running Docker

```bash
# Run the docker image, the -d tag is to run it detached
docker run -p 3001:3001 --env-file .env -d user-service

# To check the container information
docker ps

# To stop the container, use the container ID from the previous command
docker stop <container_id>
```

### Known Issues with MongoDB

The user service fails to connect to the Mongo DB server on NUS Wi-Fi. To resolve this we have to use another network.
