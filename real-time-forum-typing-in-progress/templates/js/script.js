document.addEventListener("DOMContentLoaded", function () {
  // Get the login and register buttons
  var loginBtn = document.getElementById("loginBtn");
  var registerBtn = document.getElementById("registerBtn");
  var createBtn = document.getElementById("createBtn");
  var logS = false;
  var regS = false;
  var createS = false;

  // Get the login and register forms
  var loginForm = document.getElementById("loginForm");
  var registerForm = document.getElementById("registerForm");

  // Add click event listeners to the buttons
  if (loginBtn != null) {
    loginBtn.addEventListener("click", function () {
      if (logS) {
        loginForm.style.display = "none";
        logS = false;
      } else {
        loginForm.style.display = "block";
        logS = true;
        registerForm.style.display = "none";
        regS = false;
      }
    });
  }

  if (registerBtn != null) {
    registerBtn.addEventListener("click", function () {
      if (regS) {
        registerForm.style.display = "none";
        regS = false;
      } else {
        registerForm.style.display = "block";
        regS = true;
        loginForm.style.display = "none";
        logS = false;
      }
    });
  }

  if (createBtn != null) {
    createBtn.addEventListener("click", function () {
      if (createS) {
        createForm.style.display = "none";
        createS = false;
      } else {
        createForm.style.display = "block";
        createS = true;
      }
    });
  }
});
