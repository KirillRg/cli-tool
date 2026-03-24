import http from "k6/http";

export const options = {
  stages: [
    {
      duration: "30s",
      target: 5
    },
    {
      duration: "1m",
      target: 20
    },
    {
      duration: "30s",
      target: 0
    }
  ]
};

export default function () {
  http.request("GET", "http://localtest/getting?param1=valueOfEnvParamTHREE&param3=surname", null);
  http.request("POST", "http://localtest/posting?param1=name&param3=surname", "{\n\t\"key_1\":\"value\",\n\t\"key_2\":\"valueOfEnvParamTWO\"\n}", {
    headers: {
      "Content-Type": "application/json",
      Token: "some_custom_token"
    }
  });
  http.request("DELETE", "http://localtest/deleting/valueOfEnvParamONE", null, {
    headers: {
      Token: "some_custom_token",
      New_header: "111222"
    }
  });
  http.request("PUT", "http://localtest/something", null);
}

