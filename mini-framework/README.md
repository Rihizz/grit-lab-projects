# mini-framework

This project aims to create a mini framework that abstracts the DOM routing system, implements state management, and provides event handling capabilities. The framework will be used to build a todoMVC project.

## Framework Features

The mini framework offers the following features:

1. Abstracting the DOM Routing System
2. State Management
3. Event Handling

## Documentation

The framework's documentation provides a comprehensive guide on how to use its features. Refer to the [documentation.md](documentation.md) file for detailed explanations, code examples, and insights into the framework's design choices.

## Abstracting the DOM

The framework simplifies DOM manipulation by utilizing a JavaScript object structure that represents the HTML elements. You can manipulate the DOM using methods like Virtual DOM, Data Binding, or Templating. Consider the events, children, and attributes of each element when working with the framework.

## Routing System

The framework includes a routing system that synchronizes the app's state with the URL. Changes in the URL will reflect corresponding state changes in the application.

## State Management

The framework provides state management capabilities, allowing you to handle and update the application's state. The state represents the outcome of user actions and can be accessed across multiple pages.

## Event Handling

Event handling in the framework is implemented differently from the standard `addEventListener()` method. It provides a new approach to handle user-triggered events such as scrolling, clicking, and keybindings.

## TodoMVC Project

To showcase the capabilities of the framework, a todoMVC project has been implemented using the framework. The project closely resembles the examples provided in the todoMVC website, but with the underlying implementation based on this mini framework.

## Instructions to run the Todo App

1. Clone the repository
2. run `python3 -m http.server 8000` in the terminal
3. Open `localhost:8000` in your browser
4. Add a todo item!

## Contributors

Oskar, Santeri, Stefanie, Ville and Wincent July 2023
