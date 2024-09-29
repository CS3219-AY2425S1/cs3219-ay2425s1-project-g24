This is the frontend for the question service.

## Tech Stack

- Next.js
- TypeScript
- Ant Design
- SCSS

## Getting Started

First, install the dependencies:

```bash
npm install -g pnpm

pnpm install --frozen-lockfile

# if pnpm install --frozen-lockfile fails, try running
pnpm install
```

Then, follow the `.env.example` file and create a `.env` file in the current directory. Replace the necessary values within.

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

First, run the development server:

```bash
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.
