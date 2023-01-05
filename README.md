socialnetwork

An extensible Golang social network designed for learning. The main functionality and features of the web application: authentication using JWT, the ability to add users to friends, the ability to chat with friends.

Useful functionality: 
logging of unexpected errors and events important for the application;
support for tls encryption (https);
compatibility with mongodb, the ability to configure the project.

To run the project, make sure the configuration settings in /configs are correct and start mongoDB. Scripts for running StatefulSet mongoDB on Kubernetes are provided in /deployments.
