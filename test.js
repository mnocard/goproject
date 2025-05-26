async function send() {
  const url = "http://localhost:8080/";
  const pingResponse = await fetch(url + "ping", {
    "headers": {
      "content-type": "text/plain",
      "Access-Control-Allow-Origin": "*"
    },
    "duplex": "half",
    "method": "GET",
  });
  
  const json = await pingResponse.text();
  console.log("ping: ", json);
  console.log("Response status:", pingResponse.status);

  const body = {
    value: "manu"
  }

  const adminResponse = await fetch(url + "admin", {
//   const adminResponse = await fetch(url + "admin/secrets", {
    "headers": {
      "content-type": "application/json",
    },
    body: JSON.stringify(body),
    method: "POST",
    // method: "GET",
    Authorization: 'Basic ' + btoa("admin" + ":" + "admin"),
    // Authorization: "Basic " + btoa("foo" + ":" + "bar"),
  });

  console.log("Admin response status:", adminResponse.status);
}

send();