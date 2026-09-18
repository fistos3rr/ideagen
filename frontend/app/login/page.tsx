import { LoginForm } from "@/app/login/login-form";
import { redirect } from "next/navigation";
import { headers } from 'next/headers';
import Link from 'next/link';

export default async function LoginPage() {
  const authorized = (await headers()).get('x-authorized') === 'true';

  if (authorized) {
    redirect("/");
  }

  return (
    <>
      <LoginForm />
      <p>Not registered yet? <Link href='/register'>Register</Link></p>
    </>
  );
}
