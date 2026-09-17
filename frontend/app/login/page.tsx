import { Suspense } from 'react';
import { LoginForm } from "@/app/login/login-form";
import { isAuthorized } from "@/lib/api/auth";
import { redirect } from "next/navigation";

export default function LoginPage() {
  const authorized = await isAuthorized();

  if (isAuthorized) {
    redirect("/");
  }

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <LoginForm />
    </Suspense>
  );
}
