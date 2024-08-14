// Create a function to handle route changes
function handleRouteChange(routes) {
  // console.log(routes);
  let path;
  const url = window.location.href  // get the path from the URL
  if (url.includes("#")) {
    path = url.split("#")[1]; // extract the path from the URL
  } else {
    path = url.split("/")[3];
  }

  if (path === "" || path === undefined) {
    path = "/";
  }
  // console.log(path);
  const handler = routes[path]; // Find the corresponding handler for the route

  if (handler) {
    handler(); // Invoke the handler function
  } else {
    handleNotFound(); // Handle the case when the route is not found
  }
}

// Function to handle the not found route
function handleNotFound() {
  console.log("Route not found");
}


export { handleRouteChange };