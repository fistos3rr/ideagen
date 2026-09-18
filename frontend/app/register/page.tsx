import { Suspense } from 'react';
import { RegisterForm } from "@/app/register/register-form";
import { redirect } from "next/navigation";
import { headers } from 'next/headers';

export default async function RegisterPage() {
  const authorized = (await headers()).get('x-authorized') === 'true';

  if (authorized) {
    redirect("/");
  }

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <RegisterForm />
    </Suspense>
  );
}
