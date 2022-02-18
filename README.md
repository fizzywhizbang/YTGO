# YTGO

This is a Youtube Channel monitor written in GO with a GUI writtin in QT 

## Description

I originally wrote this in Python using QT5 and developed it around the functionality of https://github.com/woefe/ytcc. The purpose of this tool
is to monitor your favorite channels, keep track of downloaded videos, and search similar content on YouTube.
Although that program worked well I wanted to develop something similar with added functionality in GO using MySQL. While that version works well for me
since I have a server at home I decided to make it more portable and have it use SQLite and perform all functions on one machine as most of the world operates.

## Getting Started

### Dependencies

* Go 1.17, QT5.3
* therecipe/qt https://github.com/therecipe/qt (follow instructions)
* and of course go mod tidy to get dependencies 


### Installing

* go mod tidy
* go mod vendor
  then
* therecipe/qt https://github.com/therecipe/qt (follow instructions)
* typically qtdeploy desktop

## License

This project is licensed under the [NAME HERE] License - see the LICENSE.md file for details

## Acknowledgments

Inspiration, code snippets, etc.
* [ytcc](https://github.com/woefe/ytcc)
* [therecipe/qt] (https://github.com/therecipe/qt)