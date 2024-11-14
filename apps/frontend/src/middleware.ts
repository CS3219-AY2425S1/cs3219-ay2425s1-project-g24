import { NextURL } from "next/dist/server/web/next-url";
import { type NextRequest, NextResponse } from "next/server";
import { jwtDecode } from "jwt-decode";

const PUBLIC_ROUTES = ["/login", "/register"];

async function isValidToken(TOKEN: string): Promise<boolean> {
  const { status } = await fetch(
    `${process.env.NEXT_PUBLIC_USER_SERVICE_URL}auth/verify-token`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${TOKEN}`,
      },
    }
  );
  return status === 200;
}

function isTokenExpired(token: string) {
  if (!token) return true;

  try {
    const decodedToken = jwtDecode(token);
    const currentTime = Date.now() / 1000; // Current time in seconds
    return decodedToken?.exp != undefined && decodedToken.exp < currentTime;
  } catch (error) {
    return true; // Return true if token is invalid
  }
}

export default async function middleware(request: NextRequest) {
  const REDIRECT_TO_LOGIN = NextResponse.redirect(
    new NextURL("/login", request.url)
  );
  const TOKEN = request.cookies.get("TOKEN");
  if (TOKEN == undefined) {
    return REDIRECT_TO_LOGIN;
  }

  if (isTokenExpired(TOKEN.value)) {
    REDIRECT_TO_LOGIN.cookies.delete("TOKEN");
    return REDIRECT_TO_LOGIN;
  }

  // FIXME: isValidToken check leads to error: not being able to access user service.
  // if (!(await isValidToken(TOKEN.value))) {
  //   REDIRECT_TO_LOGIN.cookies.delete("TOKEN");
  //   return REDIRECT_TO_LOGIN;
  // }

  return NextResponse.next();
}

export const config = {
  matcher:
    "/((?!api|_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt|login|register).*)",
  // matcher: [
  //   "/matching",
  //   "/",
  //   "/profile",
  //   "/question",
  //   "/question/.*",
  // ],
};
