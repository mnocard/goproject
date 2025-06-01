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
    "headers": {
      "content-type": "application/json",
      "authorization": 'Basic ' + btoa("admin" + ":" + "admin"),
    },
    body: JSON.stringify(body),
    method: "POST",
  });

  console.log("Admin response status:", await adminResponse.json());
  console.log("Admin response status:", adminResponse.status);
}

send();

//docker build --tag 'project' -f 'deploy/Dockerfile' .
//docker compose -f 'deploy/docker-compose.yml' up