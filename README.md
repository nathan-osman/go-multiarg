## go-multiarg

Many Go applications that I write use command-line arguments to control application behavior. The problem with command-line arguments is that they are visible to other users, eliminating the possibility of using them for sensitive data.

Environment variables and configuration files present two solutions to this problem. Unfortunately, this means writing a lot of boilerplate code to check for arguments in different places. go-multiarg exists to simplify the process.
