import { handleEvent, render, createObjectFromHTML } from "./myFramework.js";
import { handleRouteChange } from "./routing.js";
import StateManager from "./state-management.js";

// Define routes and corresponding handlers
const routes = {
  "/": handleDefaultRoute,
  "/active": handleActiveRoute,
  "/completed": handleCompletedRoute,
};

let activeFilter = "all";

const initialState = {
  todos: [],

  route: "all",

  activeFilter: "all",
};
const stateManager = new StateManager(initialState);

const styleElement = document.createElement("style");
styleElement.innerHTML = `
  .toggle-all { display: none; }
  .footer { display: none; } /* Add this rule to hide the footer */
`;
document.head.appendChild(styleElement);

const html = `
  <div>
    <section class="todoapp">
      <header class="header">
        <h1>todos</h1>
        <input class="new-todo" placeholder="What needs to be done?" autofocus>
      </header>
      <section class="main">
        <input id="toggle-all" class="toggle-all" type="checkbox">
        <label for="toggle-all">Mark all as complete</label>
        <ul class="todo-list"></ul>
      </section>
      <footer class="footer">
        <span class="todo-count">0 items left</span>
        <ul class="filters">
          <li>
            <a class="all" href="#/">All</a>
          </li>
          <li>
            <a class="act" href="#/active">Active</a>
          </li>
          <li>
            <a class="comp" href="#/completed">Completed</a>
          </li>
        </ul>
        <button class="clear-completed">Clear completed</button>
      </footer>
    </section>
    <section class="info">
      <div class="info">
        <p>Double-click to edit a todo</p>
        <p>Created by <a href="https://github.com/sbjorkman/mini-framework/">Oskar, Santeri, Stef, Ville & Wincent</a></p>
        <p>Part of <a href="http://todomvc.com">TodoMVC</a></p>
      </div>
    </section>
  </div>
`;

let renderableObject = createObjectFromHTML(html);
// console.log(stateManager.getState());
const domElement = render(renderableObject);
document.body.appendChild(domElement);

handleEvent("click", ".toggle", handleToggleTodo);
handleEvent("keydown", ".new-todo", handleAddTodo);
handleEvent("change", ".toggle-all", handleToggleAll);
handleEvent("click", ".destroy", handleDeleteTodo);
handleEvent("click", ".clear-completed", handleClearCompleted);
handleEvent("dblclick", ".todo-label", handleEditTodo);
handleEvent("click", ".filters a", handleFilter);
handleEvent("click", ".todo-list", handleTodoListClick);

// Function to render todos on page load
function renderTodos() {
  const todoList = document.querySelector(".todo-list");
  const state = stateManager.getState();
  todoList.innerHTML = "";

  state.todos.forEach((todo, index) => {
    const todoItem = createTodoItem(todo.text, index);

    if (todo.active) {
      todoItem.classList.add("active");
      //add check mark to the todo item
      const toggleCheckbox = todoItem.querySelector(".toggle");
      toggleCheckbox.checked = false;
    } else {
      todoItem.classList.add("completed");
      //add check mark to the todo item
      const toggleCheckbox = todoItem.querySelector(".toggle");
      toggleCheckbox.checked = true;
      //add clearcompleted button
      const clearCompletedButton = document.querySelector(".clear-completed");
      clearCompletedButton.style.display = "block";
    }
    todoList.appendChild(todoItem);
  });
}

renderTodos(); // Render todos on page load
updateToggleAllVisibility();

// Add a new todo item to the list
function handleAddTodo(event) {
  const newTodoInput = document.querySelector(".new-todo");
  const todoText = newTodoInput.value.trim();
  let state = stateManager.getState();

  if (event.type === "keydown" && event.key === "Enter") {
    if (todoText !== "") {
      //create state object for the new todo item
      const newTodo = {
        text: todoText,
        active: true,
        editing: false,
        dataID: state.todos.length,
      };
      //create an updated state object
      const updatedState = {
        todos: [...state.todos, newTodo],
      };
      //update the state
      stateManager.setState(updatedState);

      const todoList = document.querySelector(".todo-list");
      //call getstate to get the most up to date state with proper data ID for the new todo item
      state = stateManager.getState();
      const todoItem = createTodoItem(todoText, state.todos.length - 1);
      todoList.appendChild(todoItem);
      //add the class active to the new todo item
      todoItem.classList.add("active");
      newTodoInput.value = "";
      updateTaskCount();
      updateToggleAllVisibility();
      handleRouteChange(routes);
    }
    //print the state
    // console.log(state);
  }

  const toggleAll = document.querySelector(".toggle-all");
  toggleAll.style.display = "block"; // Always show the "toggle all" checkbox when a new item is added
}

// function to mark all items as completed
function handleToggleAll() {
  const toggleAllCheckbox = document.querySelector("#toggle-all");
  const todoList = document.querySelector(".todo-list");
  const todoItems = todoList.querySelectorAll("li");
  let state = stateManager.getState();

  todoItems.forEach((todoItem) => {
    if (toggleAllCheckbox.checked) {
      todoItem.classList.add("completed");
      todoItem.classList.remove("active");
    } else {
      todoItem.classList.remove("completed");
      todoItem.classList.add("active");
    }

    //remove it from the active list
    const currentFilter = document.querySelector(".filters .selected");
    if (currentFilter.textContent === "Active" && toggleAllCheckbox.checked) {
      todoItem.style.display = "none";
    } else if (
      currentFilter.textContent === "Active" &&
      !toggleAllCheckbox.checked
    ) {
      todoItem.style.display = "block";
    } else if (
      currentFilter.textContent === "Completed" &&
      toggleAllCheckbox.checked
    ) {
      todoItem.style.display = "block";
    } else if (
      currentFilter.textContent === "Completed" &&
      !toggleAllCheckbox.checked
    ) {
      todoItem.style.display = "none";
    } else if (
      currentFilter.textContent === "All" &&
      toggleAllCheckbox.checked
    ) {
      todoItem.style.display = "block";
    } else if (
      currentFilter.textContent === "All" &&
      !toggleAllCheckbox.checked
    ) {
      todoItem.style.display = "block";
    }

    let state = stateManager.getState();

    //update the state variable and save it to local storage
    const updatedState = {
      todos: state.todos.map((todo, index) => {
        if (index === parseInt(todoItem.dataset.id)) {
          return { ...todo, active: !toggleAllCheckbox.checked };
        }
        return todo;
      }),
    };

    stateManager.setState(updatedState);

    const toggleCheckbox = todoItem.querySelector(".toggle");
    toggleCheckbox.checked = toggleAllCheckbox.checked;
  });

  updateTaskCount();
  updateFooterVisibility();
  checkClearCompletedVisibility();
  updateToggleAllVisibility();
}

///function to handle mark one todo as completed
function handleToggleTodo(event) {
  const todoItem = event.target.closest("li");
  const isCompleted = event.target.checked;
  const state = stateManager.getState();

  if (isCompleted) {
    todoItem.classList.add("completed");
    todoItem.classList.remove("active");
    //remove it from the active list
    const currentFilter = document.querySelector(".filters .selected");
    if (currentFilter.textContent === "Active") {
      todoItem.style.display = "none";
    } else {
      todoItem.style.display = "block";
    }

    //update the state variable and save it to local storage
    const updatedState = {
      todos: state.todos.map((todo, index) => {
        if (index === parseInt(todoItem.dataset.id)) {
          return { ...todo, active: false };
        }
        return todo;
      }),
    };
    stateManager.setState(updatedState);
  } else {
    todoItem.classList.remove("completed");
    todoItem.classList.add("active");
    const currentFilter = document.querySelector(".filters .selected");
    if (currentFilter.textContent === "Completed") {
      todoItem.style.display = "none";
    } else {
      todoItem.style.display = "block";
    }
    const updatedState = {
      todos: state.todos.map((todo, index) => {
        if (index === parseInt(todoItem.dataset.id)) {
          return { ...todo, active: true };
        }
        return todo;
      }),
    };
    stateManager.setState(updatedState);
  }

  updateTaskCount();
  checkClearCompletedVisibility();
  updateToggleAllVisibility();
}

function handleTodoListClick() {
  checkClearCompletedVisibility(); // Check visibility of "Clear completed" button
}

//delete one todo item with the x button
function handleDeleteTodo(event) {
  const todoItem = event.target.closest("li");
  let state = stateManager.getState();
  // console.log("state before delete", state);

  //update the state variable and save it to local storage
  // console.log(todoItem.dataset.id);
  const updatedState = {
    todos: state.todos.filter(
      (todo, index) => index !== parseInt(todoItem.dataset.id)
    ),
  };
  stateManager.setState(updatedState);
  //update the data-id of the todo items
  // console.log("state after delete", updatedState);
  todoItem.remove();
  const todoList = document.querySelector(".todo-list");
  const todoItems = todoList.querySelectorAll("li");

  todoItems.forEach((todoItem, index) => {
    // console.log(todoItem.dataset.id);
    todoItem.dataset.id = index;
  });

  updateTaskCount();
  updateToggleAllVisibility();
  //function so show the remove all button
  checkClearCompletedVisibility();
  // console.log("handleDeleteTodo");
}

//clear all completed todos with the "Clear completed" button
function handleClearCompleted() {
  const todoList = document.querySelector(".todo-list");
  const completedItems = todoList.querySelectorAll(".completed");
  const state = stateManager.getState();

  //update the state variable and save it to local storage
  const updatedState = {
    todos: state.todos.filter((todo) => todo.active),
  };
  stateManager.setState(updatedState);

  completedItems.forEach((completedItem) => {
    todoList.removeChild(completedItem);
  });

  updateFooterVisibility();
  updateTaskCount();
  updateToggleAllVisibility();

  const clearCompletedButton = document.querySelector(".clear-completed");
  clearCompletedButton.style.display = "none";
}

//make the footer visible when a todo item is added
function updateFooterVisibility() {
  const todoList = document.querySelector(".todo-list");
  const footer = document.querySelector(".footer");
  const info = document.querySelector(".info");

  if (todoList.children.length === 0) {
    footer.style.display = "none";
    info.style.display = "block";
  } else {
    footer.style.display = "block";
    //display info
    info.style.display = "block";
  }
}

//Check if the todo items are active or completed
function checkClearCompletedVisibility() {
  const todoList = document.querySelector(".todo-list");
  const completedItems = todoList.querySelectorAll(".completed");
  const clearCompletedButton = document.querySelector(".clear-completed");

  if (completedItems.length > 0 || todoList.classList.contains("selecting")) {
    clearCompletedButton.style.display = "block";
  } else {
    clearCompletedButton.style.display = "none";
  }
}

//make the toggle all checkbox visible when a todo item is added and unmarked as all tasks are marked as completed
function updateToggleAllVisibility() {
  const todoList = document.querySelector(".todo-list");
  const toggleAll = document.querySelector(".toggle-all");
  const toggleAllLabel = document.querySelector("label[for='toggle-all']");

  if (todoList.children.length === 0) {
    toggleAll.style.display = "none";
    toggleAllLabel.style.display = "none";
  } else {
    toggleAll.style.display = "block";
    toggleAllLabel.style.display = "block";
  }
  //when not all todos are completed, the toggle all checkbox should be unchecked
  const toggleAllCheckbox = document.querySelector("#toggle-all");
  toggleAllCheckbox.checked = false;

  //when all todos are checked, the toggle all checkbox should be checked
  const todoItems = todoList.querySelectorAll("li");
  let allChecked = true;
  todoItems.forEach((todoItem) => {
    const toggleCheckbox = todoItem.querySelector(".toggle");
    if (!toggleCheckbox.checked) {
      allChecked = false;
    }
  });
  if (allChecked) {
    toggleAllCheckbox.checked = true;
  }
  //else if one clicks uncheck, the toggle all checkbox should be unchecked
  else {
    toggleAllCheckbox.checked = false;
  }
}

// edit the todo by double clicking on the todo item, esc and enter to confirm
function handleEditTodo(event) {
  let editInProgress = document.querySelector(".editing");
  if (editInProgress) {
    return;
  }
  const todoItem = event.target.closest("li");
  todoItem.classList.add("editing");
  // to the todo-label add the style none
  const todoLabel = todoItem.querySelector(".todo-label");
  todoLabel.style.display = "none";
  const editInput = todoItem.querySelector(".edit");

  editInput.value = event.target.textContent;
  editInput.focus();
  const state = stateManager.getState();

  const handleKeydown = (event, targetElement) => {
    if (event.key === "Enter") {
      const updatedState = {
        todos: state.todos.map((todo, index) => {
          if (index === parseInt(todoItem.dataset.id)) {
            return { ...todo, text: editInput.value };
          }
          return todo;
        }),
      };
      stateManager.setState(updatedState);

      handleUpdateTodo(event, targetElement.closest("li"));
      todoLabel.style.display = "block";
    }
    if (event.key === "Escape") {
      todoLabel.style.display = "block";
      handleUpdateTodo(event, targetElement.closest("li"));
    }
  };

  // Add event listener to document to handle clicks outside of the todo item
  const handleClickOutside = (event) => {
    if (!event.target.closest(".editing")) {
      todoLabel.style.display = "block";
      handleUpdateTodo(event, todoItem);
      document.removeEventListener("click", handleClickOutside);
      document.removeEventListener("keydown", handleKeydown);
    }
  };

  // Add event listeners to handle keydown and clicks outside of the todo item
  document.addEventListener("click", handleClickOutside);
  document.addEventListener("keydown", (event) =>
    handleKeydown(event, event.target)
  );
}

//update the todo item after editing with confirm with enter or esc
function handleUpdateTodo(event, todoItem) {
  const editInput = todoItem.querySelector(".edit");
  const todoLabel = todoItem.querySelector(".todo-label");
  const newText = editInput.value.trim();

  if (event.type === "keydown" && event.key === "Enter") {
    if (newText === "") {
      todoItem.remove();
    } else {
      todoLabel.textContent = newText;
    }
  }
  todoItem.classList.remove("editing");
}

//filter the todos by all, active and completed
function handleFilter(event) {
  const filters = document.querySelectorAll(".filters a");
  filters.forEach((filter) => filter.classList.remove("selected"));

  const currentFilter = event.target;
  currentFilter.classList.add("selected");

  const todoList = document.querySelector(".todo-list");
  const todoItems = todoList.querySelectorAll("li");

  todoItems.forEach((todoItem) => {
    const isCompleted = todoItem.classList.contains("completed");

    switch (currentFilter.textContent) {
      case "All":
        stateManager.setState({ activeFilter: "all" });
        todoItem.style.display = "block";
        break;
      case "Active":
        stateManager.setState({ activeFilter: "active" });
        if (!isCompleted) {
          todoItem.style.display = "block";
        } else {
          todoItem.style.display = "none";
        }
        break;
      case "Completed":
        stateManager.setState({ activeFilter: "completed" });
        if (isCompleted) {
          todoItem.style.display = "block";
        } else {
          todoItem.style.display = "none";
        }
        break;
    }
  });
}

//Create a todoItem and add it to the todo list
function createTodoItem(text, index) {
  const todoItem = document.createElement("li");
  todoItem.classList.add("todo-item");
  todoItem.setAttribute("data-id", index);

  const checkbox = document.createElement("input");
  checkbox.type = "checkbox";
  checkbox.classList.add("toggle");
  todoItem.appendChild(checkbox);

  const label = document.createElement("label");
  label.classList.add("todo-label");
  label.textContent = text;
  todoItem.appendChild(label);

  const deleteButton = document.createElement("button");
  deleteButton.classList.add("destroy");
  todoItem.appendChild(deleteButton);

  const editInput = document.createElement("input");
  editInput.type = "text";
  editInput.classList.add("edit");
  todoItem.appendChild(editInput);

  return todoItem;
}

//count the tasks left to do
function updateTaskCount() {
  const todoList = document.querySelector(".todo-list");
  const todoItems = todoList.querySelectorAll("li");
  const countElement = document.querySelector(".todo-count");
  const toggleAll = document.querySelector(".toggle-all");
  const footer = document.querySelector(".footer");

  const activeItems = Array.from(todoItems).filter(
    (todoItem) => !todoItem.classList.contains("completed")
  );

  countElement.textContent = `${activeItems.length} item${
    activeItems.length !== 1 ? "s" : ""
  } left`;

  if (todoItems.length === 0) {
    toggleAll.style.display = "none"; // Hide the "Mark all as complete" checkbox
    footer.style.display = "none"; // Hide the footer
  } else {
    toggleAll.style.display = "block";
    footer.style.display = "block"; // Show the footer
  }
}

handleEvent("hashchange", handleRouteChange(routes));

function handleActiveRoute() {
  const todoList = document.querySelector(".todo-list");
  const todoItems = todoList.querySelectorAll("li");
  const act = document.querySelector(".act");
  stateManager.setState({ activeFilter: "active" });

  activeFilter = "active";
  act.classList.add("selected");

  todoItems.forEach((todoItem) => {
    const isCompleted = todoItem.classList.contains("completed");
    const isActive = !isCompleted;

    if (activeFilter === "active") {
      todoItem.style.display = isActive ? "block" : "none";
    } else {
      todoItem.style.display = "none";
    }
  });
  updateTaskCount();
}

function handleCompletedRoute() {
  const todoList = document.querySelector(".todo-list");
  const todoItems = todoList.querySelectorAll("li");
  const comp = document.querySelector(".comp");

  activeFilter = "completed";
  comp.classList.add("selected");
  stateManager.setState({ activeFilter: "completed" });
  // console.log("im in completed route");

  todoItems.forEach((todoItem) => {
    const isCompleted = todoItem.classList.contains("completed");

    if (activeFilter === "completed") {
      todoItem.style.display = isCompleted ? "block" : "none";
    } else {
      todoItem.style.display = "none";
    }
  });
  updateTaskCount();
}

function handleDefaultRoute() {
  const todoList = document.querySelector(".todo-list");
  const todoItems = todoList.querySelectorAll("li");
  const all = document.querySelector(".all");
  // console.log(todoItems);
  stateManager.setState({ activeFilter: "all" });

  activeFilter = "all";
  all.classList.add("selected");

  todoItems.forEach((todoItem) => {
    todoItem.style.display = "block";
  });

  updateTaskCount();
  updateToggleAllVisibility(); // Move this line here

  const toggleAll = document.querySelector(".toggle-all");
  const footer = document.querySelector(".footer");

  if (todoItems.length === 0) {
    toggleAll.style.display = "none";
    footer.style.display = "none";
  } else {
    toggleAll.style.display = "block";
    footer.style.display = "block";
  }
}

handleRouteChange(routes);
updateTaskCount();
// console.log(stateManager.getState());
