# Read Imports from every Golang Project

This program takes one CLI argument which is the file path to the project you want to see all the imports for. 
I've included an example project. It will output the results of the program to a json file called imports.json 


### Prerequisites

You need to have Golang installed on your machine

```
git clone https://github.com/jmoral1943/excerise
```

## Getting Started


### Running the program

After you clone the repo make sure you are in the dir for the project and run 

```
go run main.go <filepath to the golang project relative to main.go>
```

Ex.

```
go run main.go ./paperspace-project
```

After you done running that command you can see the results in imports.json

## Program execution 

The program is broken down in a few steps

### Command line Argument 

The first part of the program takes in a CLI argument to the file path for which the project should run in.

### Recursively traverses through the dir/project folder

The listFiles function traverses any directory and searches for any golang files no matter how deep the 
folders get.

### Concurrently runs command to list all imports in a file 

```
go list -f {{.Imports}}
```

To instead of going line by line it gives me the imports for a file then I have 
splice of imports to have a JSON Object.


### Transfers the data to an JSON Object

The toJSON function converts the name of the file and the imports into an slice of maps.

### Creates and writes to a json file 

After the listFiles function is done then writeToFile runs where it will get the slice and write it to a file.