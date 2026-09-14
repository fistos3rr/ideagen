'use client';

import { useRouter } from 'next/navigation';

export function LogoutButton() {
  const router = useRouter();

  async function logout() {
    await fetch('/api/auth/logout', { method: 'POST' });
    router.push('/login');
    router.refresh();
  }

  return (
    <button
      onClick={logout}
      className="bg-red-600 text-white rounded px-3 py-2"
    >
      Выйти
    </button>
  );
}
