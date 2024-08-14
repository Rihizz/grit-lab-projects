const postComment = async (comment) => {
  if (!comment.post_id) {
    throw new Error("post_id is required");
  }

  const response = await fetch("http://localhost:6969/commenting", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(comment),
  });

  if (!response.ok) {
    const message = `An error has occurred: ${response.status}`;
    throw new Error(message);
  }

  return response.json();
};

// example usage
const comment = {
  post_id: 1,
  user_id: "1",
  content: "Hello World!",
  image: "base64 encoded image",
  created_at: "2020-01-01 00:00:00",
};

postComment(comment)
  .then((data) => console.log(data))
  .catch((error) => console.error(error));
