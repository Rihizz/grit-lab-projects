# DOM Routing System

## Creating an Element

To create an element, you can use the `document.createElement()` method. Here's an example:

```javascript
const styleElement = document.createElement("style");
styleElement.innerHTML = `.clear-completed { display: none; }`;
document.head.appendChild(styleElement);
```

In this example, a `<style>` element is created and appended to the <head> of the document. The `innerHTML` property is used to set the CSS styles for the element.

## Creating an Event

To create a event, you can use the `handleEvent()` function. Here's an example:

```javascript
function handleEvent(eventType, selector, callback) {
  document.addEventListener(eventType, (event) => {
    const targetElement = event.target.closest(selector);
    if (targetElement) {
      callback(event, targetElement);
    }
  });
}

// Usage example
handleEvent("click", ".todo-list", function (event) {
  // Event handler logic
});
```

In this example, the handleEvent function is defined to handle events of a specific type (eventType) on elements matching a specific selector (selector). The provided callback function (callback) is executed when the event occurs on a matching element.

## Nesting Elements

To nest elements, you can define a virtual DOM structure and use the `render()` and `appendChild()`` methods. Here's an example

```javascript
const virtualDOM = {
  tag: "section",
  attrs: {
    class: "todoapp",
  },
  children: [
    {
      tag: "header",
      attrs: {
        class: "header",
      },
      children: [
        // Nested elements
        {
          tag: "h1",
          attrs: {},
          children: ["todos"],
        },
        {
          tag: "input",
          attrs: {
            class: "new-todo",
            placeholder: "What needs to be done?",
            autofocus: true,
          },
          children: [],
        },
      ],
    },
    // More nested elements...
  ],
};

const app = document.getElementById("root");
const rootElement = render(virtualDOM);
app.appendChild(rootElement);
```

In this example, a virtual DOM structure is defined using an object (`virtualDOM`). The `render()` function is called to convert the virtual DOM into actual DOM elements, and then the root element is appended to the ` app`` element using  `appendChild()``.

## Creating the DOM object structure

To help with creation of the renderable objects the user can use the function createObjectFromHTML which takes html as a string and forms the object it needs to run render.

Here is an example:

```javascript
let html = `
  <div>
    <h1> Hello World </h1>
    <p> lorem ipsum </p>
  <div>
`;
let renderableObject = createObjectFromHTML(html);
const domElement = render(renderableObject);
document.body.appendChild(domElement);
```

## Adding Attributes to an Element

To add attributes to an element, you can use the `setAttribute()` method. Here's an example:

```javascript
const domElement = document.createElement(element.tag);

// Set attributes
for (const [attr, value] of Object.entries(element.attrs)) {
  domElement.setAttribute(attr, value);
}
```

In this example, the `setAttribute() method is called within a loop to set each attribute and its corresponding value on the `domElement`.

## State Management with the StateManager Class

This class handles state management tasks such as initializing, updating, getting, and clearing the application state. The state is stored in localStorage under the key "appState", and as a property of the StateManager object itself.

First, import the StateManager class into your application:
```
import StateManager from "./state-management.js"
```
Then, create an instance of the StateManager class. This will create an application state from a provided initial state or from an existing state in localStorage if one is found:
````javascript
const initialState = {
  todos: [],
  route: "all",
  activeFilter: "all",
};
const stateManager = new StateManager(initialState);
````
This method initializes the state of the application. If there's an existing state saved in localStorage under the "appState" key, that state will be used. If not, the provided initialState will be used.

`setState(newState)`
This method merges the current state with the provided newState object. After the state is updated, it is saved to localStorage:
```javascript
stateManager.setState({
  todos: ['Go shopping'],
  activeFilter: 'all'
});
````
`getState()`
This method returns the current state of the application:

`let currentState = stateManager.getState();`

`clearState()`
This method clears the current state of the application both from the StateManager instance and from localStorage:

`stateManager.clearState();`

Remember, the state maintained by the StateManager will persist across sessions due to its use of localStorage. This allows your application to maintain a consistent user experience across sessions and navigation events.