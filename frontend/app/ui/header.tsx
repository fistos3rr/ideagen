import Link from 'next/link';

export default async function Header({ authorized }: { authorized: boolean }) {
  return (
    <header className="header">
      <nav>
        <Link href="/">Main</Link>
        {authorized ? (
          <Link href="/me">My profile</Link>
        ):(
          <Link href="/login">Login</Link>
        )}
      </nav>
    </header>
  );
}
