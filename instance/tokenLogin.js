// Code taken from https://gist.github.com/m-Phoenix852/b47fffb0fd579bc210420cedbda30b61
// Complete credits to original creator (https://github.com/m-Phoenix852)

let token = "asdfgh";

function login(token) {
    setInterval(() => {
      document.body.appendChild(document.createElement `iframe`).contentWindow.localStorage.token = `"${token}"`
    }, 50);
    setTimeout(() => {
      location.reload();
    }, 2500);
  }

login(token);