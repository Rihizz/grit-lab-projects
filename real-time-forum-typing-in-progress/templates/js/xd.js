function updateShit() {
  var usersElement = document.querySelector(".users");

  // Fetch the list of online users from the server
  fetch("/api/online-users")
    .then(function (response) {
      return response.json();
    })
    .then(function (users) {
      if (users == null || users.length == 0) {
        // If no users are online, add the 'offline' class to all user elements
        usersElement.querySelectorAll(".user").forEach(function (userElement) {
          userElement.classList.add("offline");
        });
      } else {
        // Loop through each user element and set the online/offline status
        usersElement.querySelectorAll(".user").forEach(function (userElement) {
          var username = userElement.querySelector("a").textContent;
          if (users.includes(username)) {
            userElement.classList.remove("offline");
            userElement.classList.add("online");
          } else {
            userElement.classList.add("offline");
            userElement.classList.remove("online");
          }
        });
      }
    })
    .catch(function (error) {
      console.log("Error fetching online users: " + error);
    });
}
