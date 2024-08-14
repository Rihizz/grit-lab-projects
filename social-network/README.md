# social-network

In this project we created a social-network. The goal is to develop a web application with various features and functionalities commonly found in social networking platforms. The project includes both frontend and backend development using Go, JavaScript and the library React.

## Description

The social-network have the following features:

- Followers: Users are be able to follow and unfollow other users.
- Profile: Each user have a profile displaying their information, posts, and followers/following users. Profiles can be either public or private.
- Posts: Users are able to create posts and comment on existing posts. Posts can include text, images, or GIFs. Posts can have different privacy settings (public, private, or shared with selected followers).
- Groups: Users can create groups with titles and descriptions, invite other users, and accept or decline group invitations. In the group you can live-chat with the other memebers, create events and posts. Users can also be invited to the group by the creator or request to join groups, which can be accepted or declined by the group creator.
- Notifications: Users receive notifications for various events such as following requests, group invitations, and group events.
- Chats: Users can send private messages to users they are following. Group members can participate in group chats.
- Authentication: Users can register and log in using their email, password, and other optional information like avatar, nickname, and about me section. Sessions and cookies should be used for user authentication.

## Technologies

### Frontend

Frontend development involves creating the user interface and user experience of the social network using HTML, CSS, and JavaScript. We have choosen the library React.

### Backend

The backend is responsible for processing incoming requests, handling business logic and interacting with the database. It consists of the following components:

1. Server: The server receives requests and runs the application that handles those requests.
2. App: The backend application contains all the logic for processing requests, retrieving data from the database, and sending responses. It includes middleware functions that execute between receiving requests and sending responses.
3. Database: The database, implemented using SQLite, stores and organizes the data for the social network.

## Install and run

1. Clone the repository
2. Run `npm install` in the terminal
3. Split the terminal in two, one for the server and one for the application
4. Run the server in folder: backend, by typing `go run .` in the terminal
5. Run the application in folder: social-network, by typing `npm start` in the terminal

## Contributors

Oskar, Santeri, Stefanie, Ville and Wincent -
Grit:lab June 2023
