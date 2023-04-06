import Airtable from "airtable";
import { Client } from "@notionhq/client";

// ############################## contactus api ##############################
// const base = new Airtable({
//   endpointUrl: "https://api.airtable.com",
//   apiKey: process.env.NEXT_PUBLIC_AIRTABLE_API_KEY,
// }).base("appAR7Cxhflh7YVe9");

// export default base;

// ############################## ml tools ##############################

export async function fetchMLTools() {
  const notionSecret = process.env.NOTION_SECRET;
  const notionDatabaseId = process.env.NOTION_DATABASE_ID;

  const notion = new Client({ auth: notionSecret });
  if (!notionSecret || !notionDatabaseId) {
    return null;
  }

  const res = await notion.databases.query({
    database_id: notionDatabaseId,
    sorts: [
      {
        property: "Name",
        direction: "ascending",
      },
    ],
  });
  return res.results;
}

// ###########################################################################
const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL;

const makeUrl = (path: string): string => `${backendUrl}${path}`;
// ###########################################################################

// ############################ authentication api ###########################

export async function fetchSignIn(
  email: string,
  password: string,
  username?: string
) {
  const url = makeUrl(`user/login`);
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  }).then((res) => res.json());
  return res;
}

export async function fetchSignUp(
  name: string,
  username: string,
  email: string,
  password: string,
  bio?: string,
  avatar?: string
) {
  const url = makeUrl(`user/signup`);
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application / json",
    },
    body: JSON.stringify({
      name,
      handle: username,
      email,
      password,
      bio: "",
      avatar: "",
    }),
  }).then((res) => res.json());
  return res;
}

export async function fetchVerifyEmail(token: string | undefined) {
  const url = makeUrl(`user/verify-email`);
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application / json",
    },
    body: JSON.stringify({
      token: token,
    }),
  }).then((res) => res.json());
  return res;
}

export async function fetchVerifySession(
  accessToken: string,
  sessionId: string
) {
  const url = makeUrl(`user/verify-session`);
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application / json",
      Authorization: `Bearer ${accessToken}`,
    },
    body: JSON.stringify({
      session_id: sessionId,
    }),
  }).then((res) => res.json());
  return res;
}

export async function fetchForgotPassword(email: string) {
  const url = makeUrl(`user/forgot-password`);
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application / json",
    },
    body: JSON.stringify({
      email,
    }),
  }).then((res) => res.json());
  return res;
}

export async function fetchVerifyResetPassword(accessToken: string) {
  const url = makeUrl(`user/verify-reset-password`);
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application / json",
    },
    body: JSON.stringify({
      token: accessToken,
    }),
  }).then((res) => res.json());
  return res;
}

export async function fetchResetPassword(
  new_password: string,
  old_password: string,
  accessToken: string | undefined
) {
  const url = makeUrl(`user/reset-password`);
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application / json",
    },
    body: JSON.stringify({
      new_password: new_password,
      old_password: old_password,
      token: accessToken,
    }),
  }).then((res) => res.json());
  return res;
}

// ############################# user details api ############################

export async function fetchUserSettings(accessToken: string) {
  const url = makeUrl(`user/profile`);
  const res = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Accept: "application / json",
      Authorization: `Bearer ${accessToken}`,
    },
  }).then((res) => (res.ok ? res.json() : null));
  return res?.data;
}

export async function fetchPublicProfile(username: string) {
  const url = makeUrl(`user/profile/${username}`);
  const res = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Accept: "application / json",
    },
  }).then((res) => res.json());
  return res.data;
}

// ######################### update user details api #############################

export async function updateProfile(
  name: string,
  bio: string,
  accessToken: string,
  avatar?: string
) {
  const url = makeUrl(`user/profile`);
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application / json",
      Authorization: `Bearer ${accessToken}`,
    },
    body: JSON.stringify({
      name,
      bio,
      avatar: "",
    }),
  }).then((res) => res.json());
  return res;
}
