


// Navbar
function hideIconBar(){
  let iconBar = document.getElementById("iconBar");
  let navigation = document.getElementById("navigation");
  iconBar.setAttribute("style", "display:none;");
  navigation.classList.remove("hide");

}

function showIconBar(){
  let iconBar = document.getElementById("iconBar");
  let navigation = document.getElementById("navigation");
  iconBar.setAttribute("style", "display:block;");
  navigation.classList.add("hide");
}


// comment
// function showComment(){
//   let commentArea = document.getElementById("comment-area");
//   commentArea.setAttribute("style", "display:block;");
// }

// reply
function showReply(){
  let replyArea = document.getElementById("reply-area");
  replyArea.setAttribute("style", "display:block;");

}

const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');

signUpButton.addEventListener('click', () => {
container.classList.add("right-panel-active");
});

signInButton.addEventListener('click', () => {
container.classList.remove("right-panel-active");
});

// epic navbar

// function openNav() {
//     document.getElementById("mySidenav").style.width = "250px";
//   }

//   function closeNav() {
//     document.getElementById("mySidenav").style.width = "0";
//   }


// Open and close sidenav
function openNav() {
  document.getElementById("mySidenav").style.width = "250px";
  document.getElementById("main").style.marginLeft = "250px";
}

function closeNav() {
  document.getElementById("mySidenav").style.width = "0";
  document.getElementById("main").style.marginLeft = "0";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}
