import { Suspense } from 'react';
import { LoginForm } from "@/app/login/login-form";

export default function LoginPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <LoginForm />
    </Suspense>
  );
}
