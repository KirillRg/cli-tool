import http from "k6/http";

export const options = {
  vus: 1,
  duration: "10s"
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

