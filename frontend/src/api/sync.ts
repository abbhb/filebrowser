import { fetchURL } from "./utils";

export function getUserSync() {
  return fetchURL(`/api/sync/user`, { method: "GET" });
}

export function postUserSync(formdata: any) {
  return fetchURL(`/api/sync/user`, {
    method: "POST",
    body: formdata,
  });
}
