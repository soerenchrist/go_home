function createSensor(body, deviceId, res, rej) {
  fetch("/api/v1/devices/" + deviceId + "/sensors", {
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
            rej(data.message);
          })
          .catch(() => {
            rej("Error creating sensor");
          });
      }
    })
    .catch(() => {
      rej("Error creating sensor");
    });
}

function createDevice(body, res, rej) {
  fetch("/api/v1/devices", {
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
            rej(data.message);
          })
          .catch(() => {
            rej("Error creating device");
          });
      }
    })
    .catch(() => {
      rej("Error creating device");
    });
}
