"use client";

import { getToken } from "./login-store";

const USER_SERVICE_URL = process.env.NEXT_PUBLIC_USER_SERVICE_URL;

export interface User {
  id: string;
  username: string;
  email: string;
  isAdmin: boolean;
  createdAt?: string;
}

export interface UpdateUserRequestType {
  username?: string;
  password?: string;
  email?: string;
}

export interface UpdateUserResponseType {
  message: string;
  data: User;
}

export interface VerifyTokenResponseType {
  message: string;
  data: User;
}

type UserResponseData = {
  id: string;
  username: string;
  email: string;
  isAdmin: boolean;
  createdAt: string;
};

function withDefaultHeaders(opts: RequestInit): RequestInit {
  return {
    ...opts,
    headers: {
      ...(opts.headers ?? {}),
      "Content-Type": "application/json",
    },
  };
}

export async function createUser(
  username: string,
  email: string,
  password: string
): Promise<void> {
  const opts = withDefaultHeaders({
    method: "POST",
    body: JSON.stringify({
      username,
      email,
      password,
    }),
  });
  const res = await fetch(`${USER_SERVICE_URL}users`, opts);

  if (!res.ok) {
    throw new Error(
      `User service responded with ${res.status}: ${await res.text()}`
    );
  }
}

export async function loginUser(
  email: string,
  password: string
): Promise<string> {
  const opts = withDefaultHeaders({
    method: "POST",
    body: JSON.stringify({
      email,
      password,
    }),
  });
  const res = await fetch(`${USER_SERVICE_URL}auth/login`, opts);

  if (!res.ok) {
    throw new Error(
      `User service responded with ${res.status}: ${await res.text()}`
    );
  }

  const ret: { data: { accessToken: string } } = await res.json();
  return ret.data.accessToken;
}

export const UpdateUser = async (
  user: UpdateUserRequestType,
  id: string
): Promise<UpdateUserResponseType> => {
  const JWT_TOKEN = getToken();
  const response = await fetch(
    `${process.env.NEXT_PUBLIC_USER_SERVICE_URL}users/${id}`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${JWT_TOKEN}`,
      },
      body: JSON.stringify(user),
    }
  );
  if (response.status === 200) {
    return response.json();
  } else if (response.status === 409) {
    const errorResponse = await response.json();
    throw new Error(errorResponse.message);
  } else {
    throw new Error(
      `Error updating user: ${response.status} ${response.statusText}`
    );
  }
};

export const ValidateUser = async (): Promise<VerifyTokenResponseType> => {
  const JWT_TOKEN = getToken();
  const response = await fetch(
    `${process.env.NEXT_PUBLIC_USER_SERVICE_URL}auth/verify-token`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${JWT_TOKEN}`,
      },
    }
  );

  if (response.status === 200) {
    return response.json();
  } else {
    throw new Error(
      `Error verifying token: ${response.status} ${response.statusText}`
    );
  }
};
