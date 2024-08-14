document.addEventListener("DOMContentLoaded", () => {
  createForm.addEventListener("submit", handlePostFormSubmit);
  async function getPosts() {
    try {
      const response = await fetch(`/api/posts`);
      const posts = await response.json();

      if (posts == null || posts.length == 0) {
        return;
      }

      const postsContainer = document.querySelector(".post-container");
      postsContainer.innerHTML = generatePostsHTML(posts);
    } catch (error) {
      console.error(error);
    }
  }

  function generatePostsHTML(posts) {
    return posts
      .map(
        (post) => `
            <div class="post" id="post-${post.ID}" data-category="${post.Category}">
              <h2>${post.Title}</h2>
              <p class="content">${post.Content}</p>
              <p>Posted by ${post.Author} on ${post.Date} Category: 
                <a href="#" class="category-link">${post.Category}</a>
              </p>
              <a href="#" class="comment-link">Comments</a>
              <div class="comment-form-container"></div>
              <div class="comments"></div>
            </div>
          `
      )
      .join("");
  }

  function onCategoryLinkClick(event) {
    event.preventDefault();
    // Get the category of the clicked link
    const clickedCategory = event.target.innerText;

    // Get all posts
    const posts = document.querySelectorAll(".post");

    // Loop through each post and hide/show based on category
    posts.forEach((post) => {
      const postCategory = post.getAttribute("data-category");
      if (postCategory === clickedCategory) {
        post.style.display = "block";
      } else {
        post.style.display = "none";
      }
    });

    // Check if reset button already exists
    if (document.querySelector(".reset-button")) {
      return;
    }

    // Create a reset button
    const resetButton = document.createElement("button");
    resetButton.innerText = "Reset";
    resetButton.classList.add("reset-button");
    resetButton.addEventListener("click", () => {
      // Show all posts again
      posts.forEach((post) => {
        post.style.display = "block";
      });
      // Remove the reset button
      resetButton.remove();
    });

    // Append the reset button at the bottom of the body
    const postsDiv = document.querySelector(".post-container");
    postsDiv.appendChild(resetButton);
  }

  async function getComments(postID) {
    getUser().then(function (currentUsername) {
      if (currentUsername !== false) {
        const commentFormHTML = `
            <form id="comment-form" class="comment-form" action="/cum" method="POST">
                <input type="hidden" name="post_id" value="${postID}">
                <input type="hidden" name="author" value="${currentUsername}">
                <textarea name="content" placeholder="Add a comment..." required maxlength="140"></textarea>
                <button type="submit">Add Comment</button>
            </form>
        `;
        const commentFormContainer = document.querySelector(
          `#post-${postID} .comment-form-container`
        );
        commentFormContainer.innerHTML = commentFormHTML;

        const commentForm = document.querySelector(
          `#post-${postID} #comment-form`
        );
        commentForm.addEventListener("submit", handleCommentFormSubmit);
      }
      fetch(`/api/comments?post_id=${postID}`)
        .then((response) => response.json())
        .then((comments) => {
          if (comments == null || comments.length == 0) {
            return;
          }
          // Build the HTML for the comments and add it to the page
          const commentsHTML = comments
            .map(
              (comment) => `
                    <div class="comment">
                        <p class="content">${comment.Content}</p>
                        <p class="author">${comment.Author} on ${comment.Date}</p>
                    </div>
                `
            )
            .join("");
          const commentsContainer = document.querySelector(
            `#post-${postID} .comments`
          );
          commentsContainer.innerHTML = commentsHTML;
        })
        .catch((error) => {
          console.error(error);
        });
    });
  }

  function handleCommentFormSubmit(event) {
    if (event.target && event.target.id === "comment-form") {
      event.preventDefault();
      const form = event.target;
      const formData = new URLSearchParams(new FormData(form));

      const postID = formData.get("post_id");

      fetch("/cum", {
        method: "POST",
        body: formData,
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.text();
        })
        .then((data) => {
          console.log(data);
          const commentTextarea = form.querySelector("textarea");
          commentTextarea.value = "";
          const commentsContainer = document.querySelector(
            `#post-${postID} .comments`
          );
          commentsContainer.innerHTML = "";
          getComments(postID);
        })
        .catch((error) => {
          const commentTextarea = form.querySelector("textarea");
          commentTextarea.value = "";
          const commentsContainer = document.querySelector(
            `#post-${postID} .comments`
          );
          commentsContainer.innerHTML = "";
          alert(
            "The comment cannot be empty or have < or > or start or end with a space or be longer than 140 characters."
          );
          getComments(postID);
          console.error(error);
        });
    }
  }

  function onCommentLinkClick(event) {
    event.preventDefault();
    const postID = event.target.closest(".post").id.substring(5);

    const commentsContainer = document.querySelector(
      `#post-${postID} .comments`
    );
    const commentFormContainer = document.querySelector(
      `#post-${postID} .comment-form-container`
    );

    if (
      commentsContainer.innerHTML !== "" ||
      commentFormContainer.innerHTML !== ""
    ) {
      commentsContainer.innerHTML = "";
      commentFormContainer.innerHTML = "";
    } else {
      getComments(postID);
    }
  }

  function handlePostFormSubmit(event) {
    if (event.target && event.target.id === "createForm") {
      event.preventDefault();
      const form = event.target;
      const formData = new URLSearchParams(new FormData(form));

      fetch("/createPost", {
        method: "POST",
        body: formData,
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.text();
        })
        .then((data) => {
          const titleInput = form.querySelector("input[name=title]");
          titleInput.value = "";
          const contentTextarea = form.querySelector("textarea[name=content]");
          contentTextarea.value = "";
          const categoryInput = form.querySelector("input[name=category]");
          categoryInput.value = "";
          const createPostForm = document.querySelector("#createForm");
          createPostForm.style.display = "none";
          getPosts();
        })
        .catch((error) => {
          const titleInput = form.querySelector("input[name=title]");
          titleInput.value = "";
          const contentTextarea = form.querySelector("textarea[name=content]");
          contentTextarea.value = "";
          const categoryInput = form.querySelector("input[name=category]");
          categoryInput.value = "";
          alert(
            "The title/post/category cannot be empty or have < or > or start or end with a space"
          );
          console.error(error);
        });
    }
  }

  function onPostsContainerClick(event) {
    if (event.target.classList.contains("category-link")) {
      onCategoryLinkClick(event);
    } else if (event.target.classList.contains("comment-link")) {
      onCommentLinkClick(event);
    }
  }

  const postsContainer = document.querySelector(".post-container");
  postsContainer.addEventListener("click", onPostsContainerClick);

  getPosts();
});
