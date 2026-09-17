import { Suspense } from 'react';
import { LoginForm } from "@/app/login/login-form";
import { redirect } from "next/navigation";
import { headers } from 'next/headers';

export default async function LoginPage() {
  const authorized = (await headers()).get('x-authorized') === 'true';

  if (authorized) {
    redirect("/");
  }

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <LoginForm />
    </Suspense>
  );
}
