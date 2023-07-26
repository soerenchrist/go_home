function createSensor(body, deviceId, res, rej) {
  post(`/api/v1/devices/${deviceId}/sensors`, body, res, rej);
}

function createCommand(body, deviceId, res, rej) {
  post(`/api/v1/devices/${deviceId}/commands`, body, res, rej);
}

function createRule(body, res, rej) {
  post("/api/v1/rules", body, res, rej);
}

function createDevice(body, res, rej) {
  post("/api/v1/devices", body, res, rej);
}

function post(url, body, res, rej) {
  fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  })
    .then((response) => {
      if (response.status == 201) {
        res();
      } else {
        response
          .json()
          .then((data) => {
            rej(data.error);
          })
          .catch(() => {
            rej("Unknown error");
          });
      }
    })
    .catch(() => {
      rej("Unknown error");
    });
}

function showError(message) {
  let error_message = document.getElementById("error-message");
  error_message.innerHTML = message;
  error_message.classList.remove("is-hidden");
}
function hideError() {
  let error_message = document.getElementById("error-message");
  error_message.classList.add("is-hidden");
}
