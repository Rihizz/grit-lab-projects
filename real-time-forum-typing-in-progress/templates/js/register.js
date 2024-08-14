document.addEventListener("DOMContentLoaded", () => {
  const registrationForm = document.getElementById("registerForm");
  if (!registrationForm) {
    return;
  }
  const errorMessageDiv = document.getElementById("error-message");

  registrationForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    // Clear the error message if any
    errorMessageDiv.innerHTML = "";

    const formData = new FormData(registrationForm);
    const userData = Object.fromEntries(formData);

    try {
      const response = await fetch("/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(userData),
      });

      if (response.ok) {
        // Redirect to the desired page after successful registration
        window.location.href = "/";
      } else {
        // Show the error message
        const errorText = await response.text();
        errorMessageDiv.innerHTML = errorText;
      }
    } catch (error) {
      console.error("Error during registration:", error);
      errorMessageDiv.innerHTML = "An error occurred. Please try again.";
    }
  });
});
