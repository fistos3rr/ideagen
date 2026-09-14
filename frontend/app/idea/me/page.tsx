import { requireSession } from '@/lib/auth/session';
import { LogoutButton } from './logout-button';

export default async function MePage() {
  const { user } = await requireSession();

  return (
    <div className="max-w-lg mx-auto p-8 space-y-4">
      <h1 className="text-2xl font-bold">Профиль</h1>
      <p><strong>ID:</strong> {user.id}</p>
      <p><strong>Email:</strong> {user.email}</p>
      <p><strong>Роль:</strong> {user.role}</p>
      <p><strong>Создан:</strong> {new Date(user.created_at).toLocaleString('ru')}</p>

      <LogoutButton />
    </div>
  );
}
